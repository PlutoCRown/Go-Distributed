package grades

func init() {
	students =  []Student {
		{
			ID:1,
			FirstName: "A",
			LastName: "AC",
			Grades: []Grade{
				{
					Title: "Quiz1",
					Type: GradeQuiz,
					Score: 85,
				},
				{
					Title: "X",
					Type: GradeExam,
					Score: 80,
				},
				{
					Title: "Test",
					Type: GradeTest,
					Score: 90,
				},
			},
		},
	}
}