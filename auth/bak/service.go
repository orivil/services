// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package bak

import (
	limiter2 "github.com/orivil/limiter"
	"github.com/orivil/service"
	"github.com/orivil/services/captcha/email"
	"github.com/orivil/services/captcha/image"
	"github.com/orivil/services/limiter"
	"github.com/orivil/services/session"
)

type Service struct {
	storageService StorageService
	sessionService *session.Service
	imageService   *image_captcha.Service
	emailService   *email_captcha.Service
	limiterService *limiter.Service
	self           service.Provider
}

func NewService(
	storageService StorageService,
	sessionService *session.Service,
	imageCaptchaService *image_captcha.Service,
	emailCaptchaService *email_captcha.Service,
	limiterService *limiter.Service,
) *Service {
	s := &Service{
		storageService: storageService,
		sessionService: sessionService,
		imageService:   imageCaptchaService,
		emailService:   emailCaptchaService,
		limiterService: limiterService,
		self:           nil,
	}
	s.self = s
	return s
}

func (s *Service) New(ctn *service.Container) (value interface{}, err error) {
	var (
		storage           Storage
		sessionDispatcher *session.Dispatcher
		imageDispatcher   *image_captcha.Dispatcher
		emailDispatcher   *email_captcha.Dispatcher
		tlimiter          *limiter2.TimesLimiter
	)
	storage, err = s.storageService.Get(ctn)
	if err != nil {
		return nil, err
	}
	if s.sessionService != nil {
		sessionDispatcher, err = s.sessionService.Get(ctn)
		if err != nil {
			return nil, err
		}
	}
	if s.imageService != nil {
		imageDispatcher, err = s.imageService.Get(ctn)
		if err != nil {
			return nil, err
		}
	}
	if s.emailService != nil {
		emailDispatcher, err = s.emailService.Get(ctn)
		if err != nil {
			return nil, err
		}
	}
	if s.limiterService != nil {
		tlimiter, err = s.limiterService.Get(ctn)
		if err != nil {
			return nil, err
		}
	}
	return &Dispatcher{
		Store:                  storage,
		SessionDispatcher:      sessionDispatcher,
		ImageCaptchaDispatcher: imageDispatcher,
		EmailCaptchaDispatcher: emailDispatcher,
		Limiter:                tlimiter,
	}, nil
}

func (s *Service) Get(ctn *service.Container) (*Dispatcher, error) {
	dis, err := ctn.Get(&s.self)
	if err != nil {
		return nil, err
	} else {
		return dis.(*Dispatcher), nil
	}
}
