package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"
	"time"
)

type HTTPServer struct {
	db  *sql.DB
	srv *http.Server
	mux *http.ServeMux
}

func NewHTTPServer(db *sql.DB, mux *http.ServeMux) (*HTTPServer, error) {
	mux.HandleFunc("/api/something", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Not sure what you expected to find here...")
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: loggingMiddleware(mux),
	}

	s := &HTTPServer{
		db:  db,
		srv: srv,
		mux: mux,
	}
	mux.HandleFunc("GET /api/skill", s.handleSkill)
	mux.HandleFunc("POST /api/skill", s.handleSkillCreate)
	mux.HandleFunc("DELETE /api/skill/{skill_id}", s.handleSkillDelete)

	return s, nil
}

func (s *HTTPServer) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown the server
	if err := s.srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("HTTP server shutdown: %s", err)
	}
	return nil
}

func (s *HTTPServer) Start() {
	// Start the HTTP server
	log.Println("Starting HTTP server on", s.srv.Addr)
	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error running http server: %s\n", err)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		d, _ := httputil.DumpRequest(r, true)
		fmt.Println(string(d))

		// Call the next handler in the chain
		next.ServeHTTP(w, r)

		// Log after the request is completed
		duration := time.Since(start)
		log.Printf("%s - %s %s - %s - %dms\n",
			r.RemoteAddr,
			r.Method,
			r.URL,
			r.UserAgent(),
			duration.Milliseconds(),
		)
	})
}

type Skill struct {
	SkillID     int    `json:"skill_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type SkillResponse struct {
	Skills []*Skill `json:"skills"`
}

func (s *HTTPServer) handleSkill(w http.ResponseWriter, r *http.Request) {

	skills := []*Skill{}
	rows, err := s.db.Query(`
        SELECT
			skill_id, name, description
        FROM skill`)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to query skill: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var sk Skill
		err := rows.Scan(&sk.SkillID, &sk.Name, &sk.Description)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to scan skill: %v", err), http.StatusInternalServerError)
			return
		}
		skills = append(skills, &sk)
	}

	response := SkillResponse{
		Skills: skills,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal response: %v", err), http.StatusInternalServerError)
		return
	}
	w.Write(jsonResponse)
}

type SkillCreateRequest struct {
	Skill *Skill `json:"skill,omitempty"`
}

type SkillCreateResponse struct {
	Skill Skill `json:"skill"`
}

func (s *HTTPServer) handleSkillCreate(w http.ResponseWriter, r *http.Request) {

	// Unmarshal the request body into a protobuf message
	var req SkillCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Failed to unmarshal request body: %v", err), http.StatusInternalServerError)
		return
	}

	if req.Skill == nil {
		http.Error(w, "Missing skill in request body", http.StatusBadRequest)
		return
	}

	result, err := s.db.Exec(`
	INSERT INTO skill (name, description)
	VALUES (?, ?)
	`, req.Skill.Name, req.Skill.Description)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to insert skill: %v", err), http.StatusInternalServerError)
		return
	}

	// parse the result
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to lookup new skill id: %v", err), http.StatusInternalServerError)
		return
	}

	// Retrieve the inserted record
	var sk Skill
	err = s.db.QueryRow(`
	SELECT
		skill_id, name, description
	FROM skill
	WHERE skill_id = ?`, lastInsertID).Scan(
		&sk.SkillID,
		&sk.Name,
		&sk.Description,
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to lookup new skill: %v", err), http.StatusInternalServerError)
		return
	}

	resp := &SkillCreateResponse{
		Skill: sk,
	}
	jsonResponse, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal response: %v", err), http.StatusInternalServerError)
		return
	}
	w.Write(jsonResponse)
}

func (s *HTTPServer) handleSkillDelete(w http.ResponseWriter, r *http.Request) {
	sIDstr := r.PathValue("skill_id")
	sID, err := strconv.Atoi(sIDstr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse skill ID: %v", err), http.StatusBadRequest)
		return
	}

	_, err = s.db.Exec(`
	DELETE FROM skill
	WHERE skill_id = ?
	`, sID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete skill: %v", err), http.StatusInternalServerError)
		return
	}

	// TODO: dedplicate code with handleSkill
	skills := []*Skill{}
	rows, err := s.db.Query(`
        SELECT
			skill_id, name, description
        FROM skill`)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to query skill: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var sk Skill
		err := rows.Scan(&sk.SkillID, &sk.Name, &sk.Description)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to scan skill: %v", err), http.StatusInternalServerError)
			return
		}
		skills = append(skills, &sk)
	}

	response := SkillResponse{
		Skills: skills,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal response: %v", err), http.StatusInternalServerError)
		return
	}
	w.Write(jsonResponse)
}
