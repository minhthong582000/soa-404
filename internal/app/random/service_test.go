package random_test

// import (
// 	"testing"

// 	"github.com/mi"
// )

// func TestFetch(t *testing.T) {
// 	mocks.
// 	// mockRandomRepo := new(mocks.RandomRepository)
// 	// mockSeedNumber := int64(1999)

// 	// mockCases := make([]int64, 0)
// 	// mockCases = append(mockCases, mockSeedNumber)

// 	// t.Run("success", func(t *testing.T) {
// 	// 	mockRandomRepo.On("Fetch", mock.Anything, mock.AnythingOfType("string"),
// 	// 		mock.AnythingOfType("int64")).Return(mockListArtilce, "next-cursor", nil).Once()
// 	// 	mockAuthor := domain.Author{
// 	// 		ID:   1,
// 	// 		Name: "Iman Tumorang",
// 	// 	}
// 	// 	mockAuthorrepo := new(mocks.AuthorRepository)
// 	// 	mockAuthorrepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockAuthor, nil)
// 	// 	u := ucase.NewArticleUsecase(mockArticleRepo, mockAuthorrepo, time.Second*2)
// 	// 	num := int64(1)
// 	// 	cursor := "12"
// 	// 	list, nextCursor, err := u.Fetch(context.TODO(), cursor, num)
// 	// 	cursorExpected := "next-cursor"
// 	// 	assert.Equal(t, cursorExpected, nextCursor)
// 	// 	assert.NotEmpty(t, nextCursor)
// 	// 	assert.NoError(t, err)
// 	// 	assert.Len(t, list, len(mockListArtilce))

// 	// 	mockArticleRepo.AssertExpectations(t)
// 	// 	mockAuthorrepo.AssertExpectations(t)
// 	// })

// 	// t.Run("error-failed", func(t *testing.T) {
// 	// 	mockArticleRepo.On("Fetch", mock.Anything, mock.AnythingOfType("string"),
// 	// 		mock.AnythingOfType("int64")).Return(nil, "", errors.New("Unexpexted Error")).Once()

// 	// 	mockAuthorrepo := new(mocks.AuthorRepository)
// 	// 	u := ucase.NewArticleUsecase(mockArticleRepo, mockAuthorrepo, time.Second*2)
// 	// 	num := int64(1)
// 	// 	cursor := "12"
// 	// 	list, nextCursor, err := u.Fetch(context.TODO(), cursor, num)

// 	// 	assert.Empty(t, nextCursor)
// 	// 	assert.Error(t, err)
// 	// 	assert.Len(t, list, 0)
// 	// 	mockArticleRepo.AssertExpectations(t)
// 	// 	mockAuthorrepo.AssertExpectations(t)
// 	// })
// }
