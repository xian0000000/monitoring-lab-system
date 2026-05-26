package scheduler

import (
	"lab-monitor-backend/db"
	"log"
	"net"
	"os"
	"sync"
	"time"
    "github.com/go-ping/ping"
)

type ScanResult struct {
	IP        string    `json:"ip"`
	Status    string    `json:"status"`
	ScannedAt time.Time `json:"scanned_at"`
}

type sseHub struct {
	mu      sync.Mutex
	clients map[chan ScanResult]struct{}
}

var Hub = &sseHub{
	clients: make(map[chan ScanResult]struct{}),
}

func (h *sseHub) Subscribe() chan ScanResult {
	ch := make(chan ScanResult, 50)
	h.mu.Lock()
	h.clients[ch] = struct{}{}
	h.mu.Unlock()
	return ch
}

func (h *sseHub) Unsubscribe(ch chan ScanResult) {
	h.mu.Lock()
	delete(h.clients, ch)
	close(ch)
	h.mu.Unlock()
}

func (h *sseHub) broadcast(result ScanResult) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for ch := range h.clients {
		select {
		case ch <- result:
		default:
		}
	}
}

func StartPing() {
	log.Println("Ping scheduler dimulai...")
	for {
		scanSubnet()
		time.Sleep(30 * time.Second)
	}
}

func RunScanOnce() {
	scanSubnet()
}

func scanSubnet() {
	cidr := os.Getenv("NETWORK_CIDR")
	if cidr == "" {
		log.Println("NETWORK_CIDR belum diset di .env")
		return
	}

	ips, err := hostsFromCIDR(cidr)
	if err != nil {
		log.Println("CIDR tidak valid:", err)
		return
	}

	log.Printf("Mulai scan %d IP di %s\n", len(ips), cidr)

	var wg sync.WaitGroup
	sem := make(chan struct{}, 50)

	for _, ip := range ips {
		wg.Add(1)
		sem <- struct{}{}
		go func(ip string) {
			defer wg.Done()
			defer func() { <-sem }()
			pingAndBroadcast(ip)
		}(ip)
	}

	wg.Wait()
	log.Println("Scan selesai")
}

func pingAndBroadcast(ip string) {
	status := "offline"
	now := time.Now()

	pinger, err := ping.NewPinger(ip)
	if err == nil {
		pinger.Count = 1
		pinger.Timeout = 2 * time.Second
		pinger.SetPrivileged(false)

		err = pinger.Run()
		if err == nil {
			stats := pinger.Statistics()

			if stats.PacketsRecv > 0 {
				status = "online"
			}
		}
	}

	result := ScanResult{IP: ip, Status: status, ScannedAt: now}

	_, dbErr := db.DB.Exec(`
		INSERT INTO devices (name, ip_address, status, last_seen)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (ip_address)
		DO UPDATE SET status = $3, last_seen = $4
	`, ip, ip, status, now)

	if dbErr != nil {
		log.Printf("Gagal upsert %s: %v\n", ip, dbErr)
	}

	Hub.broadcast(result)
}

func hostsFromCIDR(cidr string) ([]string, error) {
	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipNet.Mask); ipNet.Contains(ip); incrementIP(ip) {
		host := ip.String()
		if host == ipNet.IP.String() {
			continue
		}
		ips = append(ips, host)
	}

	if len(ips) > 0 {
		ips = ips[:len(ips)-1]
	}

	return ips, nil
}

func incrementIP(ip net.IP) {
	for i := len(ip) - 1; i >= 0; i-- {
		ip[i]++
		if ip[i] != 0 {
			break
		}
	}
}
