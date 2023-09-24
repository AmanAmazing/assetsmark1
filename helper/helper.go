package helper 
import(
	"golang.org/x/crypto/bcrypt"
)
func HashPassword(password string) (string, error){
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password),14)
    return string(hashedPassword), err
}
func CheckPasswordsMatch(hashedPassword, enteredPassword string) bool{
    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword),[]byte(enteredPassword))
    return err == nil 
}

