package main

type Server struct {
	clientFactory ClientFactoryInterface
	dispatcher    DispatcherInterface
}

func newServer(f ClientFactoryInterface, d DispatcherInterface) *Server {
	return &Server{clientFactory: f, dispatcher: d}
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

func (s *Server) createClient(connection ConnectionInterface) ClientInterface {
	return s.clientFactory.Create(connection)
}
