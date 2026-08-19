package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/example/coldchain-fleet-dispatch-service/internal/api"
	"github.com/example/coldchain-fleet-dispatch-service/internal/config"
	"github.com/example/coldchain-fleet-dispatch-service/internal/domain"
	"github.com/example/coldchain-fleet-dispatch-service/internal/platform"
	"github.com/example/coldchain-fleet-dispatch-service/internal/service"
	"github.com/example/coldchain-fleet-dispatch-service/internal/store"
	"github.com/example/coldchain-fleet-dispatch-service/internal/worker"
)

func main() {
	cfg := config.Load()
	clk := platform.NewSystemClock()
	st := store.New()
	st.SeedFleet(defaultFleet(), defaultDrivers())

	dispatch := service.NewDispatchService(st, clk)
	temperature := service.NewTemperatureService(st, clk)
	handover := service.NewHandoverService(st, clk)
	alerts := service.NewAlertService(st, clk)
	quotes := service.NewQuoteService(st, clk)
	routes := service.NewRouteService(st, clk)
	checklists := service.NewChecklistService(st, clk)

	dispatcher := worker.NewDispatcher(st, dispatch, clk, cfg.DispatchInterval)
	alertWorker := worker.NewAlertWorker(st, alerts, clk, cfg.ExcursionWindow, cfg.AlertScanInterval)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go dispatcher.Run(ctx)
	go alertWorker.Run(ctx)

	handler := api.NewAPI(st, dispatch, temperature, handover, alerts, quotes, routes, checklists, clk).Routes()
	server := &http.Server{Addr: cfg.HTTPAddr, Handler: handler}

	go func() {
		log.Printf("coldchain fleet dispatch listening on %s", cfg.HTTPAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	_ = server.Shutdown(shutdownCtx)
}

func defaultFleet() []domain.Vehicle {
	return []domain.Vehicle{
		{ID: "veh-frozen-1", Plate: "FROZEN-01", Zone: domain.ZoneFrozen, CapacityKg: 12000, CapacityL: 32000, Status: domain.VehicleIdle},
		{ID: "veh-frozen-2", Plate: "FROZEN-02", Zone: domain.ZoneFrozen, CapacityKg: 8000, CapacityL: 22000, Status: domain.VehicleIdle},
		{ID: "veh-chilled-1", Plate: "CHILLED-01", Zone: domain.ZoneChilled, CapacityKg: 10000, CapacityL: 26000, Status: domain.VehicleIdle},
		{ID: "veh-chilled-2", Plate: "CHILLED-02", Zone: domain.ZoneChilled, CapacityKg: 9000, CapacityL: 24000, Status: domain.VehicleIdle},
		{ID: "veh-ambient-1", Plate: "AMBIENT-01", Zone: domain.ZoneAmbient, CapacityKg: 15000, CapacityL: 40000, Status: domain.VehicleIdle},
	}
}

func defaultDrivers() []domain.Driver {
	return []domain.Driver{
		{ID: "drv-frozen-1", Name: "Ada Frost", Certification: domain.ZoneFrozen, Available: true},
		{ID: "drv-frozen-2", Name: "Bo Cold", Certification: domain.ZoneFrozen, Available: true},
		{ID: "drv-chilled-1", Name: "Cy Chill", Certification: domain.ZoneChilled, Available: true},
		{ID: "drv-chilled-2", Name: "Dee Cool", Certification: domain.ZoneChilled, Available: true},
		{ID: "drv-ambient-1", Name: "Eli Dry", Certification: domain.ZoneAmbient, Available: true},
	}
}
