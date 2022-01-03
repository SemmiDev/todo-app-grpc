package random

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/SemmiDev/todo-app/model"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomTitle() string {
	return RandomString(6)
}

func RandomEmail() string {
	return fmt.Sprintf("%s@mail.com", RandomString(5))
}

func RandomDay() model.Day {
	switch randomInt(0, 6) {
	case 0:
		return model.Day_MONDAY
	case 1:
		return model.Day_TUESDAY
	case 2:
		return model.Day_WEDNESDAY
	case 3:
		return model.Day_THURSDAY
	case 4:
		return model.Day_FRIDAY
	case 5:
		return model.Day_SATURDAY
	default:
		return model.Day_SUNDAY
	}
}

func RandomDescription() string {
	return RandomString(50)
}

func RandomPriority() model.Priority {
	switch randomInt(0, 2) {
	case 0:
		return model.Priority_LOW
	case 1:
		return model.Priority_MEDIUM
	default:
		return model.Priority_HIGH
	}
}

func RandomStatus() model.Status {
	switch randomInt(0, 1) {
	case 0:
		return model.Status_IN_PROGRESS
	default:
		return model.Status_DONE
	}
}

func randomInt(min, max int) int {
	return min + rand.Int()%(max-min+1)
}
