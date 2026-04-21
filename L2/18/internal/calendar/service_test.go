package calendar

// func TestServiceAddEvent(t *testing.T) {
// 	store := NewMockEventStore(gomock.NewController(t))

// 	event := NewTestEvent("Hello, World!")

// 	store.EXPECT().Add(context.Background(), event).Return(nil)

// 	s := NewService(store)

// 	err := s.AddEvent(context.Background(), event)

// 	require.NoError(t, err)
// }

// func TestServiceUpdateEvent(t *testing.T) {
// 	store := NewMockEventStore(gomock.NewController(t))

// 	event := NewTestEvent("Hello, World!")

// 	store.EXPECT().Update(context.Background(), event).Return(nil)

// 	s := NewService(store)

// 	err := s.UpdateEvent(context.Background(), event)

// 	require.NoError(t, err)
// }
