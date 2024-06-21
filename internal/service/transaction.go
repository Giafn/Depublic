package service

import (
	"errors"

	"github.com/Giafn/Depublic/configs"
	"github.com/Giafn/Depublic/internal/entity"
	"github.com/Giafn/Depublic/internal/http/binder"
	"github.com/Giafn/Depublic/internal/repository"
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
}

type transactionService struct {
	transactionRepository repository.TransactionRepository
	db                    *gorm.DB
	cfg                   *configs.Config
}

func NewTransactionService(transactionRepository repository.TransactionRepository, db *gorm.DB, cfg *configs.Config) TransactionService {
	return &transactionService{transactionRepository: transactionRepository, db: db, cfg: cfg}
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

	newTransaction := entity.NewTransaction(eventID, userID, ticketQuantity, totalAmount, false)

	entityTickets := []entity.Ticket{}
	for _, ticket := range tickets {
		entityTicket := entity.NewTicket("", eventID.String(), ticket.BuyerName)
		entityTickets = append(entityTickets, *entityTicket)
	}
	user, err := s.transactionRepository.GetUsersById(userID)
	if err != nil {
		return nil, false, err
	}

	paymentUrl, err := s.RequestPayment(newTransaction, user)
	if err != nil {
		return nil, false, err
	}

	newTransaction.PaymentURL = paymentUrl

	transaction, err = s.transactionRepository.CreateTransaction(newTransaction, entityTickets, eventID, userID)
	if err != nil {
		return nil, false, err
	}

	// check event submission
	isMustUpload, err := s.CheckEventSubmission(eventID)
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

// request payment midtras
func (s *transactionService) RequestPayment(transaction *entity.Transaction, user *entity.User) (string, error) {
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

// check event submission
func (s *transactionService) CheckEventSubmission(eventID uuid.UUID) (bool, error) {
	event, err := s.transactionRepository.GetEventByID(eventID)
	if err != nil {
		return false, err
	}

	if !event.MustUpload {
		return false, nil
	}

	return true, nil
}
