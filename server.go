package main

type Server struct {
	idSequence IdSequenceInterface
	dispatcher DispatcherInterface
}

func newServer(s IdSequenceInterface, d DispatcherInterface) *Server {
	return &Server{idSequence: s, dispatcher: d}
}

func (s *Server) Serve(connection ConnectionInterface) {
	client := s.createClient(connection)
	defer connection.Close()

	d := s.dispatcher
	d.Subscribe(client)

	for {
		message, err := client.NextMessage()

		switch {
		case err == nil:
			d.Dispatch(message)
		case err.ConnectionError():
			d.Unsubscribe(client)
			return
			// case err.InvalidMessage(): // Just continue
		}
	}
}

func (s *Server) getNextId() uint64 {
	return s.idSequence.NextId()
}

func (s *Server) createClient(connection ConnectionInterface) ClientInterface {
	return newClient(s.getNextId(), connection)
}
