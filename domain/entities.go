package domain

type User struct {
	Id               uint64 `json:"id" gorm:"primary_key"`
	FirstName        string `json:"first_name" gorm:"not null"`
	MiddleName       string `json:"middle_name" gorm:"not null"`
	LastName         string `json:"last_name" gorm:"not null"`
	Login            string `json:"login" gorm:"not null;unique"`
	Password         string `json:"password,omitempty" gorm:"not null"`
	Email            string `json:"email" gorm:"not null;unique"`
	Occupation       string `json:"occupation"`
	IsOnaftStudent   uint64 `json:"is_onaft_student"`
	Rating           uint64 `json:"rating"`
	Role             uint64 `json:"role"`
	VerificationCode string `json:"-"`
}

type UsersTest struct {
	UserId uint64 `json:"user_id" gorm:"foreignkey:User"`
	TestId uint64 `json:"test_id" gorm:"foreignkey:ParagraphsOrTest"`
}

type Course struct {
	Id   uint64 `json:"id"   gorm:"primary_key"`
	Name string `json:"name" gorm:"not null"`
}

type Section struct {
	Id       uint64 `json:"id"   gorm:"primary_key"`
	CourseId uint64 `json:"course_id" gorm:"foreignkey:Course"`
	Name     string `json:"name" gorm:"not null"`
}

type Lesson struct {
	Id        uint64 `json:"id"   gorm:"primary_key"`
	SectionId uint64 `json:"section_id" gorm:"foreignkey:Section"`
	Name      string `json:"name" gorm:"not null"`
}

type ParagraphsOrTest struct {
	Id       uint64 `json:"id"   gorm:"primary_key"`
	LessonId uint64 `json:"lesson_id" gorm:"foreignkey:Lesson"`
	Name     string `json:"name" gorm:"not null"`
	Text     string `json:"text" gorm:"not null"`
	Points   uint64 `json:"points" gorm:"not null"`
}

type TestsAnswer struct {
	Id            uint64 `json:"id"   gorm:"primary_key"`
	TestId        uint64 `json:"test_id" gorm:"foreignkey:ParagraphsOrTest"`
	Text          string `json:"text" gorm:"not null"`
	IsRightAnswer bool   `json:"is_right_answer"`
}
