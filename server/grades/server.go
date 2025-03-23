package grades

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func RegisterHanlers(r *gin.Engine) {
	handler := new(studentsHandler)
	r.GET("/students", handler.ServeHTTP)
	r.GET("/students/", handler.ServeHTTP)

}

type studentsHandler struct{}

func (sh studentsHandler) ServeHTTP(ctx *gin.Context) {
	pathSegments := strings.Split(ctx.Request.URL.Path, "/")
	switch len(pathSegments) {
	case 2:
		sh.getAll(ctx)
	case 3:
		id, err := strconv.Atoi(pathSegments[2])
		if err != nil {
			ctx.JSON(http.StatusNotFound, nil)
			return
		}
		sh.getOne(ctx, id)
	case 4:
		id, err := strconv.Atoi(pathSegments[2])
		if err != nil {
			ctx.JSON(http.StatusNotFound, nil)
			return
		}
		sh.addGrade(ctx, id)
	default:
		ctx.JSON(http.StatusNotFound, nil)
	}
}

func (sh studentsHandler) getAll(ctx *gin.Context) {
	studentsMutex.Lock()
	defer studentsMutex.Unlock()

	data, err := sh.toJSON(students)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		log.Println(err)
		return
	}
	ctx.Writer.WriteHeader(http.StatusOK)
	ctx.Writer.Header().Add("Content-Type", "application/json")
	ctx.Writer.Write(data)
}

func (sh studentsHandler) toJSON(obj any) ([]byte, error) {
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	err := enc.Encode(obj)
	if err != nil {
		return nil, fmt.Errorf("Failed to serialize students: %q", err)
	}
	return b.Bytes(), nil
}

func (sh studentsHandler) getOne(ctx *gin.Context, id int) {
	studentsMutex.Lock()
	defer studentsMutex.Unlock()

	student, err := students.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, nil)
		log.Println(err)
		return
	}

	data, err := sh.toJSON(student)
	if err != nil {
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		log.Printf("Failed to serialize students: %q", err)
		return
	}
	ctx.Writer.WriteHeader(http.StatusOK)
	ctx.Writer.Header().Add("Content-Type", "application/json")
	ctx.Writer.Write(data)
}

func (sh studentsHandler) addGrade(ctx *gin.Context, id int) {
	studentsMutex.Lock()
	defer studentsMutex.Unlock()

	student, err := students.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, nil)
		log.Println(err)
		return
	}
	var g Grade
	dec := json.NewDecoder(ctx.Request.Body)
	err = dec.Decode(&g)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, nil)
		log.Println(err)
		return
	}
	student.Grades = append(student.Grades, g)
	ctx.Writer.WriteHeader(http.StatusCreated)
	data, err := sh.toJSON(g)
	if err != nil {
		log.Println(err)
		return
	}
	ctx.Writer.Header().Add("Content-Type", "application/json")
	ctx.Writer.Write(data)
}
