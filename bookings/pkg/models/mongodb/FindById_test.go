// ********RoostGPT********
/*
Test generated by RoostGPT for test go-unit-testing1 using AI Type  and AI Model 

ROOST_METHOD_HASH=FindByID_766b44cca3
ROOST_METHOD_SIG_HASH=FindByID_285a0c8aeb

Scenario 1: Test FindByID with valid ID

Details:
    Description: This test is meant to check if the FindByID function can successfully retrieve a booking from the database using a valid ObjectID.
Execution:
    Arrange: Mock the mongo.Collection and its FindOne method to return a valid booking when the provided ID matches the mock ID.
    Act: Call the FindByID function with the mock ID.
    Assert: Assert that the returned booking matches the mock booking and that no error is returned.
Validation:
    The assertion checks if the booking returned by the function matches the mock booking, which validates the function's ability to retrieve and return a booking from the database. This test is important as it verifies the function's basic functionality.

Scenario 2: Test FindByID with invalid ID

Details:
    Description: This test is meant to check how the FindByID function handles an invalid ObjectID.
Execution:
    Arrange: Mock the mongo.Collection and its FindOne method to return an error when the provided ID is not a valid ObjectID.
    Act: Call the FindByID function with an invalid ID.
    Assert: Assert that an error is returned.
Validation:
    The assertion checks if the function is able to handle and return an error when provided with an invalid ID. This test is important as it ensures that the function can gracefully handle invalid input.

Scenario 3: Test FindByID with non-existent ID

Details:
    Description: This test is meant to check how the FindByID function handles a non-existent ObjectID.
Execution:
    Arrange: Mock the mongo.Collection and its FindOne method to return a mongo.ErrNoDocuments error when the provided ID does not exist in the database.
    Act: Call the FindByID function with a non-existent ID.
    Assert: Assert that a "ErrNoDocuments" error is returned.
Validation:
    The assertion checks if the function is able to handle and return a custom error when the ID provided does not exist in the database. This test is important as it ensures that the function can handle situations where the requested booking does not exist.

Scenario 4: Test FindByID with database error

Details:
    Description: This test is meant to check how the FindByID function handles a database error.
Execution:
    Arrange: Mock the mongo.Collection and its FindOne method to return a generic error when called.
    Act: Call the FindByID function with a valid ID.
    Assert: Assert that an error is returned.
Validation:
    The assertion checks if the function is able to handle and return an error when a database error occurs. This test is important as it ensures that the function can handle unexpected database errors.
*/

// ********RoostGPT********
package mongodb

import (
	"context"
	"errors"
	"testing"
	"github.com/mmorejon/microservices-docker-go-mongodb/bookings/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCollection struct {
	mock.Mock
}

func (m *MockCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	args := m.Called(ctx, filter)
	return args.Get(0).(*mongo.SingleResult)
}

func TestFindById(t *testing.T) {
	mockCollection := new(MockCollection)
	mockBookingModel := &BookingModel{C: mockCollection}
	id := "5f6f7f8f9a0b1c2d3e4f5a6b"
	oid, _ := primitive.ObjectIDFromHex(id)
	mockBooking := &models.Booking{ID: oid}
	tt := []struct {
		name    string
		id      string
		mock    func()
		want    *models.Booking
		wantErr error
	}{
		{
			name: "Test FindByID with valid ID",
			id:   id,
			mock: func() {
				mockCollection.On("FindOne", context.TODO(), bson.M{"_id": oid}).Return(mockBooking, nil)
			},
			want:    mockBooking,
			wantErr: nil,
		},
		{
			name: "Test FindByID with invalid ID",
			id:   "invalidID",
			mock: func() {},
			want:    nil,
			wantErr: primitive.ErrInvalidHex,
		},
		{
			name: "Test FindByID with non-existent ID",
			id:   "5e6f7f8f9a0b1c2d3e4f5a6c",
			mock: func() {
				mockCollection.On("FindOne", context.TODO(), bson.M{"_id": oid}).Return(nil, mongo.ErrNoDocuments)
			},
			want:    nil,
			wantErr: errors.New("ErrNoDocuments"),
		},
		{
			name: "Test FindByID with database error",
			id:   id,
			mock: func() {
				mockCollection.On("FindOne", context.TODO(), bson.M{"_id": oid}).Return(nil, errors.New("database error"))
			},
			want:    nil,
			wantErr: errors.New("database error"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()
			got, err := mockBookingModel.FindByID(tc.id)
			if err != nil {
				assert.Equal(t, tc.wantErr, err)
			} else {
				assert.Equal(t, tc.want, got)
			}
		})
	}
}
