// Copyright 2022 Dolthub, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package types

import (
	"github.com/dolthub/go-mysql-server/sql/types"
)

func ConvertTypesPointToSQLPoint(p Point) types.Point {
	return types.Point{BaseGeometry: types.BaseGeometry{Geometry: p.Geometry}}
}

func ConvertTypesLineStringToSQLLineString(l LineString) types.LineString {
	return types.LineString{BaseGeometry: types.BaseGeometry{Geometry: l.Geometry}}
}

func ConvertTypesPolygonToSQLPolygon(p Polygon) types.Polygon {
	return types.Polygon{BaseGeometry: types.BaseGeometry{Geometry: p.Geometry}}
}

func ConvertTypesMultiPointToSQLMultiPoint(p MultiPoint) types.MultiPoint {
	return types.MultiPoint{BaseGeometry: types.BaseGeometry{Geometry: p.Geometry}}
}

func ConvertTypesMultiLineStringToSQLMultiLineString(l MultiLineString) types.MultiLineString {
	return types.MultiLineString{BaseGeometry: types.BaseGeometry{Geometry: l.Geometry}}
}

func ConvertTypesMultiPolygonToSQLMultiPolygon(p MultiPolygon) types.MultiPolygon {
	return types.MultiPolygon{BaseGeometry: types.BaseGeometry{Geometry: p.Geometry}}
}

func ConvertTypesGeomCollToSQLGeomColl(g GeomColl) types.GeomColl {
	return types.GeomColl{BaseGeometry: types.BaseGeometry{Geometry: g.Geometry}}
}

func ConvertSQLPointToTypesPoint(p types.Point) Point {
	return Point{Geometry: p.Geometry}
}

func ConvertSQLLineStringToTypesLineString(l types.LineString) LineString {
	return LineString{Geometry: l.Geometry}
}

func ConvertSQLPolygonToTypesPolygon(p types.Polygon) Polygon {
	return Polygon{Geometry: p.Geometry}
}

func ConvertSQLMultiPointToTypesMultiPoint(p types.MultiPoint) MultiPoint {
	return MultiPoint{Geometry: p.Geometry}
}

func ConvertSQLMultiLineStringToTypesMultiLineString(p types.MultiLineString) MultiLineString {
	return MultiLineString{Geometry: p.Geometry}
}

func ConvertSQLMultiPolygonToTypesMultiPolygon(p types.MultiPolygon) MultiPolygon {
	return MultiPolygon{Geometry: p.Geometry}
}

func ConvertSQLGeomCollToTypesGeomColl(g types.GeomColl) GeomColl {
	return GeomColl{Geometry: g.Geometry}
}

// TODO: all methods here just defer to their SQL equivalents, and assume we always receive good data

func DeserializeEWKBHeader(buf []byte) (uint32, bool, uint32, error) {
	return types.DeserializeEWKBHeader(buf)
}

func DeserializePoint(buf []byte, isBig bool, srid uint32) types.Point {
	p, _, err := types.DeserializePoint(buf, isBig, srid)
	if err != nil {
		panic(err)
	}
	return p
}

func DeserializeLine(buf []byte, isBig bool, srid uint32) types.LineString {
	l, _, err := types.DeserializeLine(buf, isBig, srid)
	if err != nil {
		panic(err)
	}
	return l
}

func DeserializePoly(buf []byte, isBig bool, srid uint32) types.Polygon {
	p, _, err := types.DeserializePoly(buf, isBig, srid)
	if err != nil {
		panic(err)
	}
	return p
}

func DeserializeMPoint(buf []byte, isBig bool, srid uint32) types.MultiPoint {
	p, _, err := types.DeserializeMPoint(buf, isBig, srid)
	if err != nil {
		panic(err)
	}
	return p
}

func DeserializeMLine(buf []byte, isBig bool, srid uint32) types.MultiLineString {
	p, _, err := types.DeserializeMLine(buf, isBig, srid)
	if err != nil {
		panic(err)
	}
	return p
}

func DeserializeMPoly(buf []byte, isBig bool, srid uint32) types.MultiPolygon {
	p, _, err := types.DeserializeMPoly(buf, isBig, srid)
	if err != nil {
		panic(err)
	}
	return p
}

func DeserializeGeomColl(buf []byte, isBig bool, srid uint32) types.GeomColl {
	g, _, err := types.DeserializeGeomColl(buf, isBig, srid)
	if err != nil {
		panic(err)
	}
	return g
}

// TODO: noms needs results to be in types

func DeserializeTypesPoint(buf []byte, isBig bool, srid uint32) Point {
	return ConvertSQLPointToTypesPoint(DeserializePoint(buf, isBig, srid))
}

func DeserializeTypesLine(buf []byte, isBig bool, srid uint32) LineString {
	return ConvertSQLLineStringToTypesLineString(DeserializeLine(buf, isBig, srid))
}

func DeserializeTypesPoly(buf []byte, isBig bool, srid uint32) Polygon {
	return ConvertSQLPolygonToTypesPolygon(DeserializePoly(buf, isBig, srid))
}

func DeserializeTypesMPoint(buf []byte, isBig bool, srid uint32) MultiPoint {
	return ConvertSQLMultiPointToTypesMultiPoint(DeserializeMPoint(buf, isBig, srid))
}

func DeserializeTypesMLine(buf []byte, isBig bool, srid uint32) MultiLineString {
	return ConvertSQLMultiLineStringToTypesMultiLineString(DeserializeMLine(buf, isBig, srid))
}

func DeserializeTypesMPoly(buf []byte, isBig bool, srid uint32) MultiPolygon {
	return ConvertSQLMultiPolygonToTypesMultiPolygon(DeserializeMPoly(buf, isBig, srid))
}

func DeserializeTypesGeomColl(buf []byte, isBig bool, srid uint32) GeomColl {
	return ConvertSQLGeomCollToTypesGeomColl(DeserializeGeomColl(buf, isBig, srid))
}
