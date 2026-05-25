// Copyright 2021 Dolthub, Inc.
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
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/dolthub/dolt/go/store/hash"
	"github.com/twpayne/go-geos"
)

// Polygon is a Noms Value wrapper around a string.
type Polygon struct {
	Geometry *geos.Geom
}

// Value interface
func (v Polygon) Value(ctx context.Context) (Value, error) {
	return v, nil
}

func (v Polygon) Equals(other Value) bool {
	if v2, ok := other.(Polygon); ok {
		return v.Geometry.Equals(v2.Geometry)
	}
	return false
}

func (v Polygon) Less(ctx context.Context, nbf *NomsBinFormat, other LesserValuable) (bool, error) {
	if v2, ok := other.(LineString); ok {
		return v.Geometry.Area() < v2.Geometry.Area(), nil
	}
	return LineStringKind < other.Kind(), nil
}

func (v Polygon) Hash(nbf *NomsBinFormat) (hash.Hash, error) {
	return getHash(v, nbf)
}

func (v Polygon) isPrimitive() bool {
	return true
}

func (v Polygon) walkRefs(nbf *NomsBinFormat, cb RefCallback) error {
	return nil
}

func (v Polygon) typeOf() (*Type, error) {
	return PrimitiveTypeMap[PolygonKind], nil
}

func (v Polygon) Kind() NomsKind {
	return PolygonKind
}

func (v Polygon) valueReadWriter() ValueReadWriter {
	return nil
}

func (v Polygon) writeTo(w nomsWriter, nbf *NomsBinFormat) error {
	err := PolygonKind.writeTo(w, nbf)
	if err != nil {
		return err
	}

	w.writeString(string(SerializePolygon(v)))
	return nil
}

func readPolygon(nbf *NomsBinFormat, b *valueDecoder) (Polygon, error) {
	buf := []byte(b.ReadString())
	srid, _, geomType, err := DeserializeEWKBHeader(buf)
	if err != nil {
		return Polygon{}, err
	}
	if geomType != WKBPolyID {
		return Polygon{}, errors.New("not a polygon")
	}
	buf = buf[EWKBHeaderSize:]
	return DeserializeTypesPoly(buf, false, srid), nil
}

func (v Polygon) readFrom(nbf *NomsBinFormat, b *binaryNomsReader) (Value, error) {
	buf := []byte(b.ReadString())
	srid, _, geomType, err := DeserializeEWKBHeader(buf)
	if err != nil {
		return Polygon{}, err
	}
	if geomType != WKBPolyID {
		return Polygon{}, errors.New("not a polygon")
	}
	buf = buf[EWKBHeaderSize:]
	return DeserializeTypesPoly(buf, false, srid), nil
}

func (v Polygon) skip(nbf *NomsBinFormat, b *binaryNomsReader) {
	b.skipString()
}

func (v Polygon) HumanReadableString() string {
	coordinateSequence := v.Geometry.CoordSeq()
	points := make([]string, coordinateSequence.Size())

	for i := range coordinateSequence.Size() {
		points[i] = fmt.Sprintf("SRID: %d POINT(%s %s)", v.Geometry.SRID(), strconv.FormatFloat(coordinateSequence.X(i), 'g', -1, 64), strconv.FormatFloat(coordinateSequence.Y(i), 'g', -1, 64))
	}
	s := fmt.Sprintf("SRID: %d POLYGON(%s)", v.Geometry.SRID(), strings.Join(points, ","))
	return strconv.Quote(s)
}
