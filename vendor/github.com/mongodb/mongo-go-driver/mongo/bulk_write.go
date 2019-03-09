// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package mongo

import (
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/mongodb/mongo-go-driver/x/mongo/driver"
)

// WriteModel is the interface satisfied by all models for bulk writes.
type WriteModel interface {
	convertModel() driver.WriteModel
}

// InsertOneModel is the write model for insert operations.
type InsertOneModel struct {
	driver.InsertOneModel
}

// NewInsertOneModel creates a new InsertOneModel.
func NewInsertOneModel() *InsertOneModel {
	return &InsertOneModel{}
}

// Document sets the BSON document for the InsertOneModel.
func (iom *InsertOneModel) Document(doc interface{}) *InsertOneModel {
	iom.InsertOneModel.Document = doc
	return iom
}

func (iom *InsertOneModel) convertModel() driver.WriteModel {
	return iom.InsertOneModel
}

// DeleteOneModel is the write model for delete operations.
type DeleteOneModel struct {
	driver.DeleteOneModel
}

// NewDeleteOneModel creates a new DeleteOneModel.
func NewDeleteOneModel() *DeleteOneModel {
	return &DeleteOneModel{}
}

// Filter sets the filter for the DeleteOneModel.
func (dom *DeleteOneModel) Filter(filter interface{}) *DeleteOneModel {
	dom.DeleteOneModel.Filter = filter
	return dom
}

// Collation sets the collation for the DeleteOneModel.
func (dom *DeleteOneModel) Collation(collation *options.Collation) *DeleteOneModel {
	dom.DeleteOneModel.Collation = collation
	return dom
}

func (dom *DeleteOneModel) convertModel() driver.WriteModel {
	return dom.DeleteOneModel
}

// DeleteManyModel is the write model for deleteMany operations.
type DeleteManyModel struct {
	driver.DeleteManyModel
}

// NewDeleteManyModel creates a new DeleteManyModel.
func NewDeleteManyModel() *DeleteManyModel {
	return &DeleteManyModel{}
}

// Filter sets the filter for the DeleteManyModel.
func (dmm *DeleteManyModel) Filter(filter interface{}) *DeleteManyModel {
	dmm.DeleteManyModel.Filter = filter
	return dmm
}

// Collation sets the collation for the DeleteManyModel.
func (dmm *DeleteManyModel) Collation(collation *options.Collation) *DeleteManyModel {
	dmm.DeleteManyModel.Collation = collation
	return dmm
}

func (dmm *DeleteManyModel) convertModel() driver.WriteModel {
	return dmm.DeleteManyModel
}

// ReplaceOneModel is the write model for replace operations.
type ReplaceOneModel struct {
	driver.ReplaceOneModel
}

// NewReplaceOneModel creates a new ReplaceOneModel.
func NewReplaceOneModel() *ReplaceOneModel {
	return &ReplaceOneModel{}
}

// Filter sets the filter for the ReplaceOneModel.
func (rom *ReplaceOneModel) Filter(filter interface{}) *ReplaceOneModel {
	rom.ReplaceOneModel.Filter = filter
	return rom
}

// Replacement sets the replacement document for the ReplaceOneModel.
func (rom *ReplaceOneModel) Replacement(rep interface{}) *ReplaceOneModel {
	rom.ReplaceOneModel.Replacement = rep
	return rom
}

// Collation sets the collation for the ReplaceOneModel.
func (rom *ReplaceOneModel) Collation(collation *options.Collation) *ReplaceOneModel {
	rom.ReplaceOneModel.Collation = collation
	return rom
}

// Upsert specifies if a new document should be created if no document matches the query.
func (rom *ReplaceOneModel) Upsert(upsert bool) *ReplaceOneModel {
	rom.ReplaceOneModel.Upsert = upsert
	rom.ReplaceOneModel.UpsertSet = true
	return rom
}

func (rom *ReplaceOneModel) convertModel() driver.WriteModel {
	return rom.ReplaceOneModel
}

// UpdateOneModel is the write model for update operations.
type UpdateOneModel struct {
	driver.UpdateOneModel
}

// NewUpdateOneModel creates a new UpdateOneModel.
func NewUpdateOneModel() *UpdateOneModel {
	return &UpdateOneModel{}
}

// Filter sets the filter for the UpdateOneModel.
func (uom *UpdateOneModel) Filter(filter interface{}) *UpdateOneModel {
	uom.UpdateOneModel.Filter = filter
	return uom
}

// Update sets the update document for the UpdateOneModel.
func (uom *UpdateOneModel) Update(update interface{}) *UpdateOneModel {
	uom.UpdateOneModel.Update = update
	return uom
}

// ArrayFilters specifies a set of filters specifying to which array elements an update should apply.
func (uom *UpdateOneModel) ArrayFilters(filters options.ArrayFilters) *UpdateOneModel {
	uom.UpdateOneModel.ArrayFilters = filters
	uom.UpdateOneModel.ArrayFiltersSet = true
	return uom
}

// Collation sets the collation for the UpdateOneModel.
func (uom *UpdateOneModel) Collation(collation *options.Collation) *UpdateOneModel {
	uom.UpdateOneModel.Collation = collation
	return uom
}

// Upsert specifies if a new document should be created if no document matches the query.
func (uom *UpdateOneModel) Upsert(upsert bool) *UpdateOneModel {
	uom.UpdateOneModel.Upsert = upsert
	uom.UpdateOneModel.UpsertSet = true
	return uom
}

func (uom *UpdateOneModel) convertModel() driver.WriteModel {
	return uom.UpdateOneModel
}

// UpdateManyModel is the write model for updateMany operations.
type UpdateManyModel struct {
	driver.UpdateManyModel
}

// NewUpdateManyModel creates a new UpdateManyModel.
func NewUpdateManyModel() *UpdateManyModel {
	return &UpdateManyModel{}
}

// Filter sets the filter for the UpdateManyModel.
func (umm *UpdateManyModel) Filter(filter interface{}) *UpdateManyModel {
	umm.UpdateManyModel.Filter = filter
	return umm
}

// Update sets the update document for the UpdateManyModel.
func (umm *UpdateManyModel) Update(update interface{}) *UpdateManyModel {
	umm.UpdateManyModel.Update = update
	return umm
}

// ArrayFilters specifies a set of filters specifying to which array elements an update should apply.
func (umm *UpdateManyModel) ArrayFilters(filters options.ArrayFilters) *UpdateManyModel {
	umm.UpdateManyModel.ArrayFilters = filters
	umm.UpdateManyModel.ArrayFiltersSet = true
	return umm
}

// Collation sets the collation for the UpdateManyModel.
func (umm *UpdateManyModel) Collation(collation *options.Collation) *UpdateManyModel {
	umm.UpdateManyModel.Collation = collation
	return umm
}

// Upsert specifies if a new document should be created if no document matches the query.
func (umm *UpdateManyModel) Upsert(upsert bool) *UpdateManyModel {
	umm.UpdateManyModel.Upsert = upsert
	umm.UpdateManyModel.UpsertSet = true
	return umm
}

func (umm *UpdateManyModel) convertModel() driver.WriteModel {
	return umm.UpdateManyModel
}

func dispatchToMongoModel(model driver.WriteModel) WriteModel {
	switch conv := model.(type) {
	case driver.InsertOneModel:
		return &InsertOneModel{conv}
	case driver.DeleteOneModel:
		return &DeleteOneModel{conv}
	case driver.DeleteManyModel:
		return &DeleteManyModel{conv}
	case driver.ReplaceOneModel:
		return &ReplaceOneModel{conv}
	case driver.UpdateOneModel:
		return &UpdateOneModel{conv}
	case driver.UpdateManyModel:
		return &UpdateManyModel{conv}
	}

	return nil
}
