package tests

import (
	"database/sql/driver"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"message_control/internal/message"
	"message_control/internal/storage"
	"message_control/internal/storage/postgres"
	"testing"
	"time"
)

const (
	testUser1 = "userT1"
	testUser2 = "userT2"
)

var (
	testMsg = message.Message{
		From: testUser1,
		To:   testUser2,
		Text: "Hallo",
		Date: time.Now(),
		Read: true,
	}

	testError = errors.New("error")
)

func TestAddNewMessage(t *testing.T) {
	testingAddNewMessage(t, nil, true)
	testingAddNewMessage(t, testError, false)
}

func testingAddNewMessage(t *testing.T, returnError error, expectedResponse bool) {
	db, mock := newDBMock(t)
	expectedQuery := `INSERT INTO`
	result := getResult()

	mock.ExpectExec(expectedQuery).WithArgs(testMsg.From, testMsg.To, testMsg.Text, testMsg.Date, testMsg.Read).WillReturnResult(result).WillReturnError(returnError)

	var controller storage.MessageControl = postgres.Postgres{Database: db}

	response := controller.AddNewMessage(testMsg)

	if response != expectedResponse {
		t.Errorf("expected: %v, got: %v, returned err: %v", expectedResponse, response, returnError)
	}
}

func newDBMock(t *testing.T) (postgres.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatal(err)
	}
	return db, mock
}

func getResult() driver.Result {
	result := sqlmock.NewResult(0, 0)
	return result
}

type testBodyGetUsers struct {
	username      string
	toMe          testCase
	sendI         testCase
	expected      []storage.ChatUser
	expectedError bool
}

type testCase struct {
	rows          [][]driver.Value
	returnedError error
}

func TestEmptyTestGetUsersList(t *testing.T) {
	test := testBodyGetUsers{
		username: "testUser",
		toMe: testCase{
			rows:          nil,
			returnedError: nil,
		},
		sendI: testCase{
			rows:          nil,
			returnedError: nil,
		},
		expected:      nil,
		expectedError: false,
	}

	testingGetUsersList(t, test)
}

func TestSuccessfulGetUsersList(t *testing.T) {
	test := testBodyGetUsers{
		username: "testUser",
		toMe: testCase{
			rows: [][]driver.Value{
				{"N1", 10},
				{"N2", 0},
				{"N3", 1},
				{"N10", 0},
			},
			returnedError: nil,
		},
		sendI: testCase{
			rows: [][]driver.Value{
				{"N4"},
				{"N5"},
				{"N6"},
			},
			returnedError: nil,
		},
		expected: []storage.ChatUser{
			{"N1", false},
			{"N2", true},
			{"N3", false},
			{"N10", true},
			{"N4", false},
			{"N5", false},
			{"N6", false},
		},
		expectedError: false,
	}

	testingGetUsersList(t, test)
}

func TestErrorsGetUsersList(t *testing.T) {
	tests := []testBodyGetUsers{
		{
			username: "testUser",
			toMe: testCase{
				rows:          nil,
				returnedError: testError,
			},
			sendI: testCase{
				rows:          nil,
				returnedError: testError,
			},
			expected:      nil,
			expectedError: true,
		},
		{
			username: "testUser",
			toMe: testCase{
				rows:          nil,
				returnedError: nil,
			},
			sendI: testCase{
				rows:          nil,
				returnedError: testError,
			},
			expected:      nil,
			expectedError: true,
		},
		{
			username: "testUser",
			toMe: testCase{
				rows:          nil,
				returnedError: testError,
			},
			sendI: testCase{
				rows:          nil,
				returnedError: nil,
			},
			expected:      nil,
			expectedError: true,
		},
	}

	for _, test := range tests {
		testingGetUsersList(t, test)
	}
}

var (
	toMeColumns = []string{
		"from_user",
		"count_false",
	}

	sendIColumns = []string{
		"to_user",
	}
)

func testingGetUsersList(t *testing.T, test testBodyGetUsers) {
	db, mock := newDBMock(t)
	expectedQuery := `SELECT`

	// test getUsersMessagesToMe
	rowsMock := createRows(toMeColumns, test.toMe.rows)
	mock.ExpectQuery(expectedQuery).WithArgs(test.username).WillReturnRows(rowsMock).WillReturnError(test.toMe.returnedError)

	// test getUsersMessagesSendI
	rowsMock = createRows(sendIColumns, test.sendI.rows)
	mock.ExpectQuery(expectedQuery).WithArgs(test.username, test.username).WillReturnRows(rowsMock).WillReturnError(test.sendI.returnedError)

	var controller storage.MessageControl = postgres.Postgres{Database: db}
	responseUsers, err := controller.GetUsersList(test.username)

	checkError(t, test.expectedError, err)
	if !test.expectedError {
		equalsChatUser(t, test.expected, responseUsers)
	}
}

func createRows(columns []string, rows [][]driver.Value) *sqlmock.Rows {
	rowsMock := sqlmock.NewRows(columns)

	for _, row := range rows {
		rowsMock.AddRow(row...)
	}
	return rowsMock
}

func checkError(t *testing.T, expectedError bool, err error) {
	if expectedError {

		if err == nil {
			t.Fatal("expected error, got nil")
		}

	} else {

		if err != nil {
			t.Fatal("expected nil, got ", err)
		}
	}
}

func equalsChatUser(t *testing.T, expected, response []storage.ChatUser) {
	checkLen(t, expected, response)

	var (
		e, r storage.ChatUser
	)

	for i := range expected {
		e, r = expected[i], response[i]

		if e.Username != r.Username {
			t.Errorf("expected username %s, got %s", e.Username, r.Username)
		}

		if e.Read != r.Read {
			t.Errorf("expected read %v, got %v", e.Read, r.Read)
		}
	}
}

func checkLen[k []storage.ChatUser | []message.Message](t *testing.T, expected, response k) {
	lnExpected := len(expected)
	lnResponse := len(response)

	if lnExpected != lnResponse {
		t.Error("Expected: ", expected)
		t.Error("Response: ", response)
		t.Fatalf("expected len: %d, got len: %d", lnExpected, lnResponse)
	}
}

func TestEmptyGetMessagesChat(t *testing.T) {
	test := testBodyGetMessagesChat{
		username: "testUser",
		buddy:    "testBuddy",
		table: testCase{
			rows:          nil,
			returnedError: nil,
		},
		expected:      nil,
		expectedError: false,
	}
	testingGetMessagesChat(t, test)
}

var (
	getMessagesColumns = []string{"from", "to", "text", "date", "read"}
)

type testBodyGetMessagesChat struct {
	username      string
	buddy         string
	table         testCase
	expected      []message.Message
	expectedError bool
}

func testingGetMessagesChat(t *testing.T, test testBodyGetMessagesChat) {
	db, mock := newDBMock(t)

	rowsMock := createRows(getMessagesColumns, test.table.rows)
	mock.ExpectQuery("SELECT").WithArgs(test.username, test.buddy).WillReturnRows(rowsMock).WillReturnError(test.table.returnedError)

	var controller storage.MessageControl = postgres.Postgres{Database: db}

	responseUsers, err := controller.GetMessagesChat(test.username, test.buddy)

	checkError(t, test.expectedError, err)

	if !test.expectedError {
		equalsMessages(t, test.expected, responseUsers)
	}
}

func equalsMessages(t *testing.T, expected, response []message.Message) {
	checkLen(t, expected, response)

	var (
		e, r message.Message
	)

	for i := range expected {
		e, r = expected[i], response[i]

		if e.From != r.From {
			t.Errorf("expected from %s, got %s", e.From, r.From)
		}

		if e.To != r.To {
			t.Errorf("expected to %s, got %s", e.To, r.To)
		}

		if e.Text != r.Text {
			t.Errorf("expected text %s, got %s", e.Text, r.Text)
		}

		if e.Date != r.Date {
			t.Errorf("expected date %s, got %s", e.Date, r.Date)
		}

		if e.Read != r.Read {
			t.Errorf("expected read %v, got %v", e.Read, r.Read)
		}
	}
}

func TestErrorGetMessagesChat(t *testing.T) {
	test := testBodyGetMessagesChat{
		username: "testUser",
		buddy:    "testBuddy",
		table: testCase{
			rows:          nil,
			returnedError: testError,
		},
		expected:      nil,
		expectedError: true,
	}
	testingGetMessagesChat(t, test)
}
