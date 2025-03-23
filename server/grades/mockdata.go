package grades

func init() {
	students = []Student{
		{
			ID:        1,
			FirstName: "Huang",
			LastName:  "Yi",
			Grades: []Grade{
				{
					Title: "Quiz 1",
					Type:  GradeQuiz,
					Score: 85,
				},
				{
					Title: "Final exam",
					Type:  GradeExam,
					Score: 95,
				},
			},
		},
		{
			ID:        2,
			FirstName: "Peter",
			LastName:  "HH",
			Grades: []Grade{
				{
					Title: "Quiz 1",
					Type:  GradeQuiz,
					Score: 99,
				},
				{
					Title: "Final exam",
					Type:  GradeExam,
					Score: 98,
				},
			},
		},
	}
}
