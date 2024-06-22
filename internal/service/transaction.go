package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Giafn/Depublic/configs"
	"github.com/Giafn/Depublic/internal/entity"
	"github.com/Giafn/Depublic/internal/http/binder"
	"github.com/Giafn/Depublic/internal/repository"
	"github.com/Giafn/Depublic/pkg/encrypt"
	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"gorm.io/gorm"
)

type TransactionService interface {
	CreateTransaction(eventID uuid.UUID, userID uuid.UUID, tickets []binder.Ticket) (*entity.Transaction, bool, error)
	FindTransactionByID(id uuid.UUID) (*entity.Transaction, error)
	UpdateTransaction(transaction *entity.Transaction) (*entity.Transaction, error)
	FindAllTransactions() ([]entity.Transaction, error)
	DeleteTransaction(id uuid.UUID) error
	CountAmountTickets(tickets []binder.Ticket, eventID uuid.UUID) (int, error)
	GetUsersById(id uuid.UUID) (*entity.User, error)
	FindMyTransactions(userID uuid.UUID) ([]entity.Transaction, error)
	EncryptPaymentURL(paymentURL string, transactionID uuid.UUID) (string, error)
	DecryptPaymentURL(encryptedPaymentURL string) (string, error)
	CheckTicketAvailability(transactionID uuid.UUID) (bool, error)
	UpdateTicketRemaining(tickets []entity.Ticket) error
	GetTicketsByTransactionID(transactionID uuid.UUID) ([]entity.Ticket, error)
}

type transactionService struct {
	transactionRepository repository.TransactionRepository
	db                    *gorm.DB
	encryptTool           encrypt.EncryptTool
	cfg                   *configs.Config
}

func NewTransactionService(transactionRepository repository.TransactionRepository,
	db *gorm.DB,
	encryptTool encrypt.EncryptTool,
	cfg *configs.Config,
) TransactionService {
	return &transactionService{transactionRepository: transactionRepository, db: db, encryptTool: encryptTool, cfg: cfg}
}

func (s *transactionService) CreateTransaction(
	eventID uuid.UUID,
	userID uuid.UUID,
	tickets []binder.Ticket,
) (transaction *entity.Transaction, mustUpload bool, err error) {
	totalAmount, err := s.CountAmountTickets(tickets, eventID)
	if err != nil {
		return nil, false, err
	}
	ticketQuantity := len(tickets)

	for _, ticket := range tickets {
		pricing, err := s.transactionRepository.GetPricingByEventID(eventID, ticket.PricingId)
		if err != nil {
			return nil, false, err
		}

		if pricing.Remaining < 1 {
			return nil, false, errors.New("ticket not available")
		}
	}

	newTransaction := entity.NewTransaction(eventID, userID, ticketQuantity, totalAmount, false)

	entityTickets := []entity.Ticket{}
	for _, ticket := range tickets {
		entityTicket := entity.NewTicket("", eventID.String(), ticket.BuyerName, ticket.PricingId.String())
		entityTickets = append(entityTickets, *entityTicket)
	}
	user, err := s.transactionRepository.GetUsersById(userID)
	if err != nil {
		return nil, false, err
	}

	paymentUrl, err := s.requestPayment(newTransaction, user)
	if err != nil {
		return nil, false, err
	}

	newTransaction.PaymentURL = paymentUrl

	transaction, err = s.transactionRepository.CreateTransaction(newTransaction, entityTickets, eventID, userID)
	if err != nil {
		return nil, false, err
	}

	// check event submission
	isMustUpload, err := s.checkEventSubmission(eventID)
	if err != nil {
		return nil, false, err
	}

	return transaction, isMustUpload, nil
}

func (s *transactionService) FindTransactionByID(id uuid.UUID) (*entity.Transaction, error) {
	transaction, err := s.transactionRepository.FindTransactionByID(id)
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

func (s *transactionService) UpdateTransaction(transaction *entity.Transaction) (*entity.Transaction, error) {
	existingTransaction, err := s.transactionRepository.FindTransactionByID(transaction.ID)
	if err != nil {
		return nil, errors.New("transaction not found")
	}

	existingTransaction.EventID = transaction.EventID
	existingTransaction.UserID = transaction.UserID
	existingTransaction.TicketQuantity = transaction.TicketQuantity
	existingTransaction.TotalAmount = transaction.TotalAmount
	existingTransaction.IsPaid = transaction.IsPaid

	updatedTransaction, err := s.transactionRepository.UpdateTransaction(existingTransaction)
	if err != nil {
		return nil, err
	}
	return updatedTransaction, nil
}

func (s *transactionService) FindAllTransactions() ([]entity.Transaction, error) {
	transactions, err := s.transactionRepository.FindAllTransactions()
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (s *transactionService) DeleteTransaction(id uuid.UUID) error {
	return s.transactionRepository.DeleteTransaction(id)
}

func (s *transactionService) CountAmountTickets(tickets []binder.Ticket, eventID uuid.UUID) (int, error) {
	var totalAmount int
	for _, ticket := range tickets {
		pricing, err := s.transactionRepository.GetPricingByEventID(eventID, ticket.PricingId)
		if err != nil {
			return 0, err
		}
		totalAmount += pricing.Fee
	}
	return totalAmount, nil
}

func (s *transactionService) requestPayment(transaction *entity.Transaction, user *entity.User) (string, error) {
	snapClient := snap.Client{}
	serverKey := s.cfg.Midtrans.ServerKey
	snapClient.New(serverKey, midtrans.Sandbox)

	request := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  transaction.ID.String(),
			GrossAmt: int64(transaction.TotalAmount),
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: user.Profiles.FullName,
			LName: "(customer depublic)",
			Email: user.Email,
		},
	}

	snapResponse, err := snapClient.CreateTransaction(request)
	if err != nil {
		return "", err
	}

	paymentUrl := snapResponse.RedirectURL

	return paymentUrl, nil
}

func (s *transactionService) checkEventSubmission(eventID uuid.UUID) (bool, error) {
	event, err := s.transactionRepository.GetEventByID(eventID)
	if err != nil {
		return false, err
	}

	if !event.MustUpload {
		return false, nil
	}

	return true, nil
}

func (s *transactionService) GetUsersById(id uuid.UUID) (*entity.User, error) {
	user, err := s.transactionRepository.GetUsersById(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *transactionService) FindMyTransactions(userID uuid.UUID) ([]entity.Transaction, error) {
	transactions, err := s.transactionRepository.FindTransactionsByUserId(userID)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (s *transactionService) EncryptPaymentURL(paymentURL string, transactionID uuid.UUID) (string, error) {
	explodedPaymentURL := strings.Split(paymentURL, "/")

	paymentId := explodedPaymentURL[len(explodedPaymentURL)-1]

	url := fmt.Sprintf("http://%s:%s/app/api/v1/payment?pay_id=%s&transaction_id=%s", s.cfg.Host, s.cfg.Port, paymentId, transactionID.String())
	return url, nil
}

func (s *transactionService) DecryptPaymentURL(paymentId string) (string, error) {

	url := fmt.Sprintf("https://app.sandbox.midtrans.com/snap/v4/redirection/%s", paymentId)
	return url, nil
}

func (s *transactionService) CheckTicketAvailability(transactionID uuid.UUID) (bool, error) {
	transaction, err := s.transactionRepository.FindTransactionByID(transactionID)
	if err != nil {
		return false, err
	}

	if transaction.IsPaid {
		return false, errors.New("transaction already paid")
	}

	ticket, err := s.transactionRepository.FindTicketByTransactionID(transactionID)
	if err != nil {
		return false, err
	}

	for _, t := range ticket {
		pricing, err := s.transactionRepository.GetPricingById(uuid.MustParse(t.PricingID))
		if err != nil {
			return false, err
		}

		if pricing.Remaining < 1 {
			return false, errors.New("ticket not available")
		}
	}

	return true, nil
}

func (s *transactionService) UpdateTicketRemaining(tickets []entity.Ticket) error {
	for _, ticket := range tickets {
		pricing, err := s.transactionRepository.GetPricingById(uuid.MustParse(ticket.PricingID))
		if err != nil {
			return err
		}

		pricing.Remaining -= 1

		_, err = s.transactionRepository.UpdatePricingRemaining(pricing.PricingId, pricing.Remaining)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *transactionService) GetTicketsByTransactionID(transactionID uuid.UUID) ([]entity.Ticket, error) {
	tickets, err := s.transactionRepository.FindTicketByTransactionID(transactionID)
	if err != nil {
		return nil, err
	}
	return tickets, nil
}
