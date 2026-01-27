package services

import (
    "database/sql"
    "fmt"
    
    "github.com/FLASH2332/Secure-Question-Paper-Distribution-Portal/internal/acl"
    "github.com/FLASH2332/Secure-Question-Paper-Distribution-Portal/internal/models"
)

// FacultyService handles faculty-specific operations
type FacultyService struct {
    DB   *sql.DB
    User *models.User
}

// NewFacultyService creates a new faculty service
func NewFacultyService(db *sql.DB, user *models.User) *FacultyService {
    return &FacultyService{
        DB:   db,
        User: user,
    }
}

// CanUploadPaper checks if faculty can upload papers
func (s *FacultyService) CanUploadPaper() error {
    return acl.EnforcePermission(s.DB, s.User, "QuestionPaper", "create", nil)
}

// CanEncrypt checks if faculty can encrypt papers
func (s *FacultyService) CanEncrypt() error {
    return acl.EnforcePermission(s.DB, s.User, "QuestionPaper", "encrypt", nil)
}

// CanViewOwnPapers checks if faculty can view their papers
func (s *FacultyService) CanViewOwnPapers() error {
    return acl.EnforcePermission(s.DB, s.User, "QuestionPaper", "read", nil)
}

// GetMyPapers retrieves papers uploaded by this faculty
func (s *FacultyService) GetMyPapers() ([]models.QuestionPaper, error) {
    // First check permission
    if err := s.CanViewOwnPapers(); err != nil {
        return nil, err
    }
    
    query := `
        SELECT id, title, subject, faculty_id, upload_date, exam_date, status
        FROM question_papers
        WHERE faculty_id = ?
        ORDER BY upload_date DESC
    `
    
    rows, err := s.DB.Query(query, s.User.ID)
    if err != nil {
        return nil, fmt.Errorf("failed to get papers: %w", err)
    }
    defer rows.Close()
    
    var papers []models.QuestionPaper
    for rows.Next() {
        var paper models.QuestionPaper
        var examDate sql.NullTime
        
        err := rows.Scan(
            &paper.ID,
            &paper.Title,
            &paper.Subject,
            &paper.FacultyID,
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
