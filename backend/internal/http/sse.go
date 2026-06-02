package http

import (
	"context"
	"log"

	"feedsystem_video_go/internal/middleware/rabbitmq"
	"feedsystem_video_go/internal/worker"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setupSSE(r *gin.Engine, db *gorm.DB, rmq *rabbitmq.RabbitMQ) {
	if rmq != nil && rmq.Ch != nil {
		rmq.DeclareTopic("like.events", "notification.like", "like.like")
		rmq.DeclareTopic("comment.events", "notification.comment", "comment.publish")
		rmq.DeclareTopic("social.events", "notification.social", "social.follow")
	}
	sseHub := worker.NewSSEHub(db)
	notifGroup := r.Group("/notification")
	notifGroup.Use(sseHub.SSERequireAuth())
	sseHub.RegisterRoutes(r, notifGroup)

	go func() {
		if rmq != nil && rmq.Ch != nil {
			hub := sseHub
			ctx := context.Background()
			go func() {
				w := worker.NewNotificationWorker(rmq.Ch, db, "notification.like", hub)
				if err := w.Run(ctx); err != nil {
					log.Printf("notification-like worker: %v", err)
				}
			}()
			go func() {
				w := worker.NewNotificationWorker(rmq.Ch, db, "notification.comment", hub)
				if err := w.Run(ctx); err != nil {
					log.Printf("notification-comment worker: %v", err)
				}
			}()
			go func() {
				w := worker.NewNotificationWorker(rmq.Ch, db, "notification.social", hub)
				if err := w.Run(ctx); err != nil {
					log.Printf("notification-social worker: %v", err)
				}
			}()
		} else {
			log.Printf("Notification SSE disabled (MQ not available)")
		}
	}()
}
