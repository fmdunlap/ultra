package log

import (
    "errors"
    "fmt"
    "io"
    "strconv"
    "testing"
    "time"
)

type user struct {
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
}

func fakeMessages(n int) []string {
    messages := make([]string, n)
    for i := range messages {
        messages[i] = fmt.Sprintf("Test logging, but use a somewhat realistic message length. (#%v)", i)
    }
    return messages
}

func fakeVals() []any {
    return []any{
        _tenInts[0],
        _tenInts,
        _tenStrings[0],
        _tenStrings,
        _tenTimes[0],
        _tenTimes,
        _oneUser,
        _oneUser,
        _tenUsers,
        errExample,
    }
}

func getMessage(iter int) string {
    return _messages[iter%1000]
}

type users []*user

var (
    errExample = errors.New("fail")

    _messages   = fakeMessages(1000)
    _tenInts    = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
    _tenStrings = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
    _tenTimes   = []time.Time{
        time.Unix(0, 0),
        time.Unix(1, 0),
        time.Unix(2, 0),
        time.Unix(3, 0),
        time.Unix(4, 0),
        time.Unix(5, 0),
        time.Unix(6, 0),
        time.Unix(7, 0),
        time.Unix(8, 0),
        time.Unix(9, 0),
    }
    _oneUser = &user{
        Name:      "Jane Doe",
        Email:     "jane@test.com",
        CreatedAt: time.Date(1980, 1, 1, 12, 0, 0, 0, time.UTC),
    }
    _tenUsers = users{
        _oneUser,
        _oneUser,
        _oneUser,
        _oneUser,
        _oneUser,
        _oneUser,
        _oneUser,
        _oneUser,
        _oneUser,
        _oneUser,
    }
)

// Benchmark test for logging to Info
func BenchmarkLogger_Log_oneField(b *testing.B) {
    formatter, _ := NewFormatter(OutputFormatText, []Field{NewDefaultLevelField(), NewMessageField()})
    logger, _ := NewLoggerWithOptions(WithDestination(io.Discard, formatter), WithMinLevel(Info), WithAsync(false))

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        logger.Info("test")
    }
}

func BenchmarkLogger_Log_TenFields(b *testing.B) {
    intField, _ := NewIntField("int")
    intsField, _ := NewArrayField[int]("ints", func(args LogLineArgs, data int) (any, error) {
        if args.OutputFormat == OutputFormatText {
            return strconv.Itoa(data), nil
        }
        return data, nil
    })
    stringField, _ := NewStringField("string")
    stringsField, _ := NewArrayField[string]("strings", func(args LogLineArgs, data string) (any, error) {
        return data, nil
    })
    timeFIeld, _ := NewTimeField("time", "2006-01-02 15:04:05")
    timesField, _ := NewArrayField[time.Time]("times", func(args LogLineArgs, data time.Time) (any, error) {
        if args.OutputFormat == OutputFormatText {
            return data.Format("2006-01-02 15:04:05"), nil
        }
        return data, nil
    })
    userField, _ := NewObjectField[user]("user", func(args LogLineArgs, data user) (any, error) {
        if args.OutputFormat == OutputFormatText {
            return fmt.Sprintf("'%s'", data), nil
        }

        return data, nil
    })
    usersField, _ := NewArrayField[user]("users", func(args LogLineArgs, data user) (any, error) {
        if args.OutputFormat == OutputFormatText {
            return fmt.Sprintf("'%s'", data.Name), nil
        }

        return data, nil
    })

    errorField, _ := NewErrorField("error")

    formatter, _ := NewFormatter(OutputFormatText, []Field{
        intField,
        intsField,
        stringField,
        stringsField,
        timeFIeld,
        timesField,
        userField,
        userField,
        usersField,
        errorField,
    })

    logger, _ := NewLoggerWithOptions(WithDestination(io.Discard, formatter), WithMinLevel(Info), WithAsync(false))

    b.ResetTimer()
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            logger.Log(Info, fakeVals()...)
        }
    })
}
