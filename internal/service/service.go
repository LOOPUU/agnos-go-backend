package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/LOOPUU/agnos-hospital-middleware/internal/client"
	"github.com/LOOPUU/agnos-hospital-middleware/internal/model"
	"github.com/LOOPUU/agnos-hospital-middleware/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials=errors.New("invalid credentials")
type Service struct { Repo *repository.Repository; Hospital *client.HospitalClient; TokenTTL time.Duration }
func (s *Service) CreateStaff(ctx context.Context,hospital,username,password string)(model.Staff,error){hash,err:=bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost);if err!=nil{return model.Staff{},err};return s.Repo.CreateStaff(ctx,hospital,username,string(hash))}
func (s *Service) Login(ctx context.Context,hospital,username,password string)(model.Token,error){staff,hash,err:=s.Repo.LoginData(ctx,hospital,username);if err!=nil||bcrypt.CompareHashAndPassword([]byte(hash),[]byte(password))!=nil{return model.Token{},ErrInvalidCredentials};b:=make([]byte,32);if _,err=rand.Read(b);err!=nil{return model.Token{},err};raw:=hex.EncodeToString(b);expires:=time.Now().Add(s.TokenTTL);if err=s.Repo.SaveToken(ctx,staff.ID,raw,expires);err!=nil{return model.Token{},err};return model.Token{AccessToken:raw,TokenType:"Bearer",ExpiresAt:expires},nil}
func (s *Service) Authenticate(ctx context.Context,token string)(model.Staff,error){return s.Repo.Authenticate(ctx,token)}
func (s *Service) Search(ctx context.Context,staff model.Staff,f model.SearchFilters)([]model.Patient,error){id:=f.NationalID;if id==""{id=f.PassportID};if id!=""&&staff.APIBaseURL!=nil{p,err:=s.Hospital.Find(ctx,*staff.APIBaseURL,id);if err!=nil{return nil,err};if p!=nil{if err=s.Repo.UpsertPatient(ctx,staff.HospitalID,*p);err!=nil{return nil,err}}};return s.Repo.SearchPatients(ctx,staff.HospitalID,f)}

