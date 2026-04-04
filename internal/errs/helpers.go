package errs

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func BadRequestFromConstraint(constraint string) *AppError {
	switch constraint {
	case "fk_categories_parent_id":
		return BadRequest(errors.New("there is no category with such parent id"))
	case "fk_products_brand_id":
		return BadRequest(errors.New("there is no brand with such id"))
	case "fk_products_category_id":
		return BadRequest(errors.New("there is no category with such id"))
	default:
		return Internal(errors.New("constraint\"" + constraint + "\"not handled in BadRequestFromConstraint function"))
	}
}

func ConflictFromConstraint(constraint string) *AppError {
	switch constraint {
	case "brands_name_key":
		return Conflict(errors.New("brand's name already exists"))
	case "brands_slug_key":
		return Conflict(errors.New("brand's slug already exists"))
	case "categories_slug_key":
		return Conflict(errors.New("category's slug already exists"))
	case "products_slug_key":
		return Conflict(errors.New("product's slug already exists"))
	case "admin_users_email_key":
		return Conflict(errors.New("admin_user's email already exists"))
	default:
		return Internal(errors.New("constraint\"" + constraint + "\"not handled in ConflictFromConstraint function"))
	}
}

func FromPgErr(err error) *AppError {
	var pgErr *pgconn.PgError
	switch {
	case errors.As(err, &pgErr):
		switch pgErr.Code {
		// unique_violation
		case "23505":
			return ConflictFromConstraint(pgErr.ConstraintName)
		// foreign_key_violation
		case "23503":
			return BadRequestFromConstraint(pgErr.ConstraintName)
		default:
			return Internal(fmt.Errorf("unknown pgconn.PgError.Code %v: %w", pgErr.Code, pgErr))
		}
	case errors.Is(err, pgx.ErrNoRows):
		return NotFound(err)
	default:
		return Internal(fmt.Errorf("error is not processable by FromPgErr: %w", err))
	}
}
