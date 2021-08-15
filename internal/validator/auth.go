package validator

//type authValidator struct {
//	db *sqlx.DB
//}
//
//func (a authValidator) EmailIsUnique(email interface{}) error {
//	var count int64
//	if err := a.db.Get(&count, `SELECT count(1) FROM "user" WHERE email = ?`, email); err != nil {
//		return err
//	} else if count != 0 {
//		return errors.New("email already exist")
//	}
//
//	return nil
//}
//
//func NewAuthValidator(db *sqlx.DB) domain.AuthValidator {
//	return authValidator{db: db}
//}
