package services

import (
	"database/sql"
	"fmt"

	"github.com/FLASH2332/Secure-Question-Paper-Distribution-Portal/internal/acl"
	"github.com/FLASH2332/Secure-Question-Paper-Distribution-Portal/internal/models"
)

// ExamCellService handles exam cell operations
type ExamCellService struct {
	DB   *sql.DB
	User *models.User
}

// NewExamCellService creates a new exam cell service
func NewExamCellService(db *sql.DB, user *models.User) *ExamCellService {
	return &ExamCellService{
		DB:   db,
		User: user,
	}
}

// CanDecryptPaper checks if exam cell can decrypt papers
func (s *ExamCellService) CanDecryptPaper() error {
	return acl.EnforcePermission(s.DB, s.User, "QuestionPaper", "decrypt", nil)
}

// CanViewAllPapers checks if exam cell can view all papers
func (s *ExamCellService) CanViewAllPapers() error {
	return acl.EnforcePermission(s.DB, s.User, "QuestionPaper", "read", nil)
}

// CanCreateSession checks if exam cell can create exam sessions
func (s *ExamCellService) CanCreateSession() error {
	return acl.EnforcePermission(s.DB, s.User, "ExamSession", "create", nil)
}

// GetAllPapers retrieves all question papers
func (s *ExamCellService) GetAllPapers() ([]models.QuestionPaper, error) {
	// First check permission
	if err := s.CanViewAllPapers(); err != nil {
		return nil, err
	}

	query := `
        SELECT qp.id, qp.title, qp.subject, qp.faculty_id, u.username as faculty_name,
               qp.upload_date, qp.exam_date, qp.status
        FROM question_papers qp
        JOIN users u ON qp.faculty_id = u.id
        ORDER BY qp.upload_date DESC
    `

	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get papers: %w", err)
	}
	defer rows.Close()

	var papers []models.QuestionPaper
	for rows.Next() {
		var paper models.QuestionPaper
		var facultyName string
		var examDate sql.NullTime

		err := rows.Scan(
			&paper.ID,
			&paper.Title,
			&paper.Subject,
			&paper.FacultyID,
			&facultyName,
			&paper.UploadDate,
			&examDate,
			&paper.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan paper: %w", err)
		}

		if examDate.Valid {
			paper.ExamDate = examDate.Time
		}

		papers = append(papers, paper)
	}

	return papers, nil
}
