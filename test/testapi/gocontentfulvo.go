// Package testapi - DO NOT EDIT THIS FILE: Auto-generated code by https://github.com/foomo/gocontentful
package testapi

import (
	"sync"

	"github.com/foomo/contentful"
)

type CfBrand struct {
	Sys    ContentfulSys     `json:"sys,omitempty"`
	Fields CfBrandFields     `json:"fields,omitempty"`
	CC     *ContentfulClient `json:"-"`
}

// CfBrandFields is a CfNameFields VO
type CfBrandFields struct {
	CompanyName              map[string]string         `json:"companyName,omitempty"`
	RWLockCompanyName        sync.RWMutex              `json:"-"`
	Logo                     map[string]ContentTypeSys `json:"logo,omitempty"`
	RWLockLogo               sync.RWMutex              `json:"-"`
	CompanyDescription       map[string]string         `json:"companyDescription,omitempty"`
	RWLockCompanyDescription sync.RWMutex              `json:"-"`
	Website                  map[string]string         `json:"website,omitempty"`
	RWLockWebsite            sync.RWMutex              `json:"-"`
	Twitter                  map[string]string         `json:"twitter,omitempty"`
	RWLockTwitter            sync.RWMutex              `json:"-"`
	Email                    map[string]string         `json:"email,omitempty"`
	RWLockEmail              sync.RWMutex              `json:"-"`
	Phone                    map[string][]string       `json:"phone,omitempty"`
	RWLockPhone              sync.RWMutex              `json:"-"`
}

type CfBrandFieldsLogo struct {
	Entry *contentful.Asset
	Col   *contentful.Collection
}

type CfCategory struct {
	Sys    ContentfulSys     `json:"sys,omitempty"`
	Fields CfCategoryFields  `json:"fields,omitempty"`
	CC     *ContentfulClient `json:"-"`
}

// CfCategoryFields is a CfNameFields VO
type CfCategoryFields struct {
	Title                     map[string]string         `json:"title,omitempty"`
	RWLockTitle               sync.RWMutex              `json:"-"`
	Icon                      map[string]ContentTypeSys `json:"icon,omitempty"`
	RWLockIcon                sync.RWMutex              `json:"-"`
	CategoryDescription       map[string]string         `json:"categoryDescription,omitempty"`
	RWLockCategoryDescription sync.RWMutex              `json:"-"`
}

type CfCategoryFieldsIcon struct {
	Entry *contentful.Asset
	Col   *contentful.Collection
}

type CfProduct struct {
	Sys    ContentfulSys     `json:"sys,omitempty"`
	Fields CfProductFields   `json:"fields,omitempty"`
	CC     *ContentfulClient `json:"-"`
}

// CfProductFields is a CfNameFields VO
type CfProductFields struct {
	ProductName              map[string]string           `json:"productName,omitempty"`
	RWLockProductName        sync.RWMutex                `json:"-"`
	Slug                     map[string]string           `json:"slug,omitempty"`
	RWLockSlug               sync.RWMutex                `json:"-"`
	ProductDescription       map[string]string           `json:"productDescription,omitempty"`
	RWLockProductDescription sync.RWMutex                `json:"-"`
	Sizetypecolor            map[string]string           `json:"sizetypecolor,omitempty"`
	RWLockSizetypecolor      sync.RWMutex                `json:"-"`
	Image                    map[string][]ContentTypeSys `json:"image,omitempty"`
	RWLockImage              sync.RWMutex                `json:"-"`
	Tags                     map[string][]string         `json:"tags,omitempty"`
	RWLockTags               sync.RWMutex                `json:"-"`
	Categories               map[string][]ContentTypeSys `json:"categories,omitempty"`
	RWLockCategories         sync.RWMutex                `json:"-"`
	Price                    map[string]float64          `json:"price,omitempty"`
	RWLockPrice              sync.RWMutex                `json:"-"`
	Brand                    map[string]ContentTypeSys   `json:"brand,omitempty"`
	RWLockBrand              sync.RWMutex                `json:"-"`
	Quantity                 map[string]float64          `json:"quantity,omitempty"`
	RWLockQuantity           sync.RWMutex                `json:"-"`
	Sku                      map[string]string           `json:"sku,omitempty"`
	RWLockSku                sync.RWMutex                `json:"-"`
	Website                  map[string]string           `json:"website,omitempty"`
	RWLockWebsite            sync.RWMutex                `json:"-"`
}

type genericEntryNoFields struct {
	Sys ContentfulSys `json:"sys,omitempty"`
}
