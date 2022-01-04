package random

import (
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"math/rand"
	"strings"
	"time"

	"github.com/SemmiDev/todo-app/proto"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func RandomID() string {
	return uuid.NewString()
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomCreatedAt() *timestamppb.Timestamp {
	now := timestamppb.Now()
	yearAgo := now.AsTime().AddDate(-1, 0, 0)
	return timestamppb.New(time.Unix(RandomInt(yearAgo.Unix(), now.AsTime().Unix()), 0))
}

func RandomUpdatedAt() *timestamppb.Timestamp {
	now := timestamppb.Now()
	yearAgo := now.AsTime().AddDate(0, -5, 0)
	return timestamppb.New(time.Unix(RandomInt(yearAgo.Unix(), now.AsTime().Unix()), 0))
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

func RandomDay() proto.Day {
	switch randomInt(0, 6) {
	case 0:
		return proto.Day_MONDAY
	case 1:
		return proto.Day_TUESDAY
	case 2:
		return proto.Day_WEDNESDAY
	case 3:
		return proto.Day_THURSDAY
	case 4:
		return proto.Day_FRIDAY
	case 5:
		return proto.Day_SATURDAY
	default:
		return proto.Day_SUNDAY
	}
}

func RandomDescription() string {
	return RandomString(50)
}

func RandomPriority() proto.Priority {
	switch randomInt(0, 2) {
	case 0:
		return proto.Priority_LOW
	case 1:
		return proto.Priority_MEDIUM
	default:
		return proto.Priority_HIGH
	}
}

func RandomStatus() proto.Status {
	switch randomInt(0, 1) {
	case 0:
		return proto.Status_IN_PROGRESS
	default:
		return proto.Status_DONE
	}
}

func randomInt(min, max int) int {
	return min + rand.Int()%(max-min+1)
}
