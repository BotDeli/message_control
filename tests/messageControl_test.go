package tests

import (
	"database/sql/driver"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"message_control/internal/message"
	"message_control/internal/server/serverGRPC/pb"
	"message_control/internal/storage"
	"message_control/internal/storage/postgres"
	"message_control/pkg/format"
	"testing"
	"time"
)

type testCase[T any] struct {
	username      string
	friend        string
	returnedRows  [][]driver.Value
	returnedError error
	expected      []T
	expectedError bool
}

const (
	testUser1 = "userT1"
	testUser2 = "userT2"
	SELECT    = `SELECT`
	INSERT    = `INSERT INTO`
)

var (
	testError = errors.New("error")
)

func TestAddNewMessage(t *testing.T) {
	testingAddNewMessage(t, nil, true, false)
	testingAddNewMessage(t, testError, false, true)
}

func testingAddNewMessage(t *testing.T, returnError error, expectedResponse, expectedError bool) {
	db, mock := newDBMock(t)
	result := getResult()

	testMsg := message.Message{
		From: testUser1,
		To:   testUser2,
		Text: "Hallo",
		Date: time.Now(),
	}

	mock.ExpectExec(INSERT).WithArgs(testMsg.From, testMsg.To, testMsg.Text, testMsg.Date).WillReturnResult(result).WillReturnError(returnError)

	var controller storage.MessageControl = postgres.Postgres{Database: db}

	response, err := controller.AddNewMessage(testMsg)

	checkError(t, expectedError, err)

	if response != expectedResponse {
		t.Errorf("expected: %v, got: %v, returned err: %v", expectedResponse, response, returnError)
	}
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

func TestEmptyTestGetUsersList(t *testing.T) {
	test := testCase[*pb.Friend]{
		username:      "testUser",
		returnedRows:  [][]driver.Value{},
		returnedError: nil,
		expected:      []*pb.Friend{},
		expectedError: false,
	}

	testingGetFriendsList(t, test)
}

func testingGetFriendsList(t *testing.T, test testCase[*pb.Friend]) {
	db, mock := newDBMock(t)

	rowsMock := createRows([]string{"username", "date"}, test.returnedRows)
	mock.ExpectQuery(SELECT).WithArgs(test.username).WillReturnRows(rowsMock).WillReturnError(test.returnedError)

	var controller storage.MessageControl = postgres.Postgres{Database: db}
	responseUsers, err := controller.GetFriendsList(test.username)

	checkError(t, test.expectedError, err)

	equalFriends(t, test.expected, responseUsers)
}

func createRows(columns []string, rows [][]driver.Value) *sqlmock.Rows {
	rowsMock := sqlmock.NewRows(columns)

	for _, row := range rows {
		rowsMock.AddRow(row...)
	}
	return rowsMock
}

func equalFriends(t *testing.T, expected, response []*pb.Friend) {
	checkLen(t, expected, response)
	for i := range expected {

		if expected[i].Username != response[i].Username {
			t.Errorf("expected username %s, got %s", expected[i].Username, response[i].Username)
		}

		if expected[i].Date != response[i].Date {
			t.Errorf("expected date %s, got %s", expected[i].Date, response[i].Date)
		}

	}
}

func checkLen[T any](t *testing.T, expected, response []T) {
	if len(expected) != len(response) {
		t.Fatalf("expected length %d, response length %d", len(expected), len(response))
	}
}

func TestSuccessfulGetUsersList(t *testing.T) {
	testTime := time.Now()

	test := testCase[*pb.Friend]{
		username: "testUser",
		returnedRows: [][]driver.Value{
			{"user1", testTime},
			{"user2", testTime},
			{"user3", testTime},
			{"user4", testTime},
			{"user5", testTime},
			{"user6", testTime},
			{"user7", testTime},
			{"user8", testTime},
			{"user9", testTime},
			{"user10", testTime},
			{"u1", testTime},
			{"u2", testTime},
		},
		returnedError: nil,
		expected: []*pb.Friend{
			{Username: "user1", Date: format.Date(testTime)},
			{Username: "user2", Date: format.Date(testTime)},
			{Username: "user3", Date: format.Date(testTime)},
			{Username: "user4", Date: format.Date(testTime)},
			{Username: "user5", Date: format.Date(testTime)},
			{Username: "user6", Date: format.Date(testTime)},
			{Username: "user7", Date: format.Date(testTime)},
			{Username: "user8", Date: format.Date(testTime)},
			{Username: "user9", Date: format.Date(testTime)},
			{Username: "user10", Date: format.Date(testTime)},
			{Username: "u1", Date: format.Date(testTime)},
			{Username: "u2", Date: format.Date(testTime)},
		},
		expectedError: false,
	}

	testingGetFriendsList(t, test)
}

func TestErrorsGetUsersList(t *testing.T) {
	test := testCase[*pb.Friend]{
		username:      "testUser",
		returnedRows:  nil,
		returnedError: testError,
		expected:      nil,
		expectedError: true,
	}

	testingGetFriendsList(t, test)
}

func TestEmptyGetMessagesChat(t *testing.T) {
	test := testCase[*pb.ChatMessage]{
		username:      "testUser",
		friend:        "testFriend",
		returnedRows:  [][]driver.Value{},
		returnedError: nil,
		expected:      []*pb.ChatMessage{},
		expectedError: false,
	}
	testingGetMessagesChat(t, test)
}

func testingGetMessagesChat(t *testing.T, test testCase[*pb.ChatMessage]) {
	db, mock := newDBMock(t)
	rowsMock := createRows([]string{"from", "to", "text", "date"}, test.returnedRows)

	mock.ExpectQuery(SELECT).WithArgs(test.username, test.friend).WillReturnRows(rowsMock).WillReturnError(test.returnedError)

	var controller storage.MessageControl = postgres.Postgres{Database: db}

	responseUsers, err := controller.GetMessagesChat(test.username, test.friend)

	checkError(t, test.expectedError, err)

	equalChatMessage(t, test.expected, responseUsers)
}

func equalChatMessage(t *testing.T, expected, response []*pb.ChatMessage) {
	checkLen(t, expected, response)

	for i := range expected {

		if expected[i].Msg.To != response[i].Msg.To {
			t.Errorf("expected msg.to %s, got %s", expected[i].Msg.To, response[i].Msg.To)
		}

		if expected[i].Msg.From != response[i].Msg.From {
			t.Errorf("expected msg.from %s, got %s", expected[i].Msg.From, response[i].Msg.From)
		}

		if expected[i].Msg.Text != response[i].Msg.Text {
			t.Errorf("expected msg.text %s, got %s", expected[i].Msg.Text, response[i].Msg.Text)
		}

		if expected[i].Date != response[i].Date {
			t.Errorf("expected date %s, got %s", expected[i].Date, response[i].Date)
		}

	}
}

func TestErrorGetMessagesChat(t *testing.T) {
	test := testCase[*pb.ChatMessage]{
		username:      "testUser",
		friend:        "testBuddy",
		returnedRows:  nil,
		returnedError: testError,
		expected:      nil,
		expectedError: true,
	}
	testingGetMessagesChat(t, test)
}

func TestSuccessfulGetMessagesChat(t *testing.T) {
	testTime := time.Now()

	test := testCase[*pb.ChatMessage]{
		username: "testUser",
		friend:   "testBuddy",
		returnedRows: [][]driver.Value{
			{"testUser", "testBuddy", "", testTime},
			{"testBuddy", "testUser", "HelloWorld", testTime},
			{"testUser", "testBuddy", "HelloTest :3 .... \n OK! Process...\n NEW YEAR!!!", testTime},
			{"testUser", "testBuddy", "No! No! No!", testTime},
			{"testBuddy", "testUser", "ccc", testTime},
		},
		returnedError: nil,
		expected: []*pb.ChatMessage{
			{
				Msg: &pb.BodyMessage{
					From: "testUser",
					To:   "testBuddy",
					Text: "",
				},
				Date: format.Date(testTime),
			},
			{
				Msg: &pb.BodyMessage{
					From: "testBuddy",
					To:   "testUser",
					Text: "HelloWorld",
				},
				Date: format.Date(testTime),
			},
			{
				Msg: &pb.BodyMessage{
					From: "testUser",
					To:   "testBuddy",
					Text: "HelloTest :3 .... \n OK! Process...\n NEW YEAR!!!",
				},
				Date: format.Date(testTime),
			},
			{
				Msg: &pb.BodyMessage{
					From: "testUser",
					To:   "testBuddy",
					Text: "No! No! No!",
				},
				Date: format.Date(testTime),
			},
			{
				Msg: &pb.BodyMessage{
					From: "testBuddy",
					To:   "testUser",
					Text: "ccc",
				},
				Date: format.Date(testTime),
			},
		},
		expectedError: false,
	}
	testingGetMessagesChat(t, test)
}
