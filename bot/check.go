package bot

//func CheckStaleAnswers(staleQuestionCheckInterval int, staleQuestionAfter int) {
//	for {
//		select {
//		case <-time.After(time.Duration(staleQuestionCheckInterval) * time.Second):
//			ids := db.CheckLastQuestions(staleQuestionAfter)
//			msg := staleQuestionConversations(CS.SupportChat, ids)
//			if _, err := B.Api.Send(msg); err != nil {
//				log.L.Error(err)
//			}
//		}
//	}
//}
