package app

import (
	"avito/internal/config"
	"avito/internal/handler"
	"avito/internal/middleware"
	"avito/internal/repository"
	"avito/internal/usecase"
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
)

// Run запускает приложение и HTTP-сервер.
func Run() {
	service := gin.Default()
	service.Use(middleware.ServerMiddleware())

	// init config
	cfg, err := config.New()
	fmt.Println(cfg)
	if err != nil {
		panic("Error loading config" + err.Error())
	}

	// db init
	db, err := repository.Init(context.Background(), cfg)
	if err != nil {
		panic("Error initializing repository" + err.Error())
	}
	defer db.Close()

	// repository
	userRepo := repository.NewUserRepository(db.DB())
	roomRepo := repository.NewRoomRepository(db.DB())
	scheduleRepo := repository.NewScheduleRepository(db.DB())
	slotRepo := repository.NewSlotsRepository(db.DB())
	bookingRepo := repository.NewBookingRepository(db.DB())

	// usecase
	confUC := usecase.NewConferenceUseCase()
	authUC := usecase.NewAuthUseCase(userRepo, cfg.JWTSecret)
	roomUC := usecase.NewRoomUseCase(roomRepo)
	scheduleUC := usecase.NewScheduleUseCase(scheduleRepo, roomRepo)
	slotUC := usecase.NewSlotUseCase(slotRepo, scheduleRepo, roomRepo)
	bookingUC := usecase.NewBookingService(bookingRepo, slotRepo, confUC)

	// middleware
	mw := middleware.NewAuthMiddleware(authUC)

	// handlers
	h := handler.New(authUC, roomUC, scheduleUC, slotUC, bookingUC)

	// Public ручки
	service.POST("/dummyLogin", h.DummyLogin)
	service.POST("/register", h.Register)
	service.POST("/login", h.Login)
	service.GET("/_info", h.Info)

	// группируем ручки, где нужен jwt
	auth := service.Group("/")
	auth.Use(mw.Authenticate())

	// Rooms
	rooms := auth.Group("/rooms")
	rooms.GET("/list", h.ListRooms)
	rooms.POST("/create", h.CreateRoom)

	// Schedules
	rooms.POST("/:roomId/schedule/create", h.CreateSchedule)

	// Slots
	rooms.GET("/:roomId/slots/list", h.ListSlots)

	// Bookings
	bookings := auth.Group("/bookings")
	bookings.POST("/create", h.CreateBooking)
	bookings.GET("/list", h.ListBookings)
	bookings.GET("/my", h.MyBookings)
	bookings.POST("/:bookingId/cancel", h.CancelBooking)

	if err := service.Run(":8080"); err != nil {
		panic("Error starting service" + err.Error())
	}
}
