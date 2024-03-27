package exception

import (
	"github.com/gin-gonic/gin"
	"gitlab.cs.ui.ac.id/ahmadhi.prananta/execution_service/validation"
)

func ErrorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func ErrorResponseMessage(err string) gin.H {
	return gin.H{"error": err}
}

func ViolationsFieldValidation(violations []validation.BadRequest_FieldViolation) gin.H {
	return gin.H{"message": "invalid parameters", "details": violations}
}

func ViolationUniqueConstraint(errConstraint string) gin.H {
	return gin.H{"message": "unique constraint violation", "details": errConstraint}
}

func ServerErrorResponse(err error) gin.H {
	return gin.H{"message": "internal server error", "details": err.Error()}
}
