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

// LineString is a Noms Value wrapper around a string.
type LineString struct {
	Geometry *geos.Geom
}

// Value interface
func (v LineString) Value(ctx context.Context) (Value, error) {
	return v, nil
}

func (v LineString) Equals(other Value) bool {
	if v2, ok := other.(LineString); ok {
		return v.Geometry.Equals(v2.Geometry)
	}
	return false
}

func (v LineString) Less(ctx context.Context, nbf *NomsBinFormat, other LesserValuable) (bool, error) {
	if v2, ok := other.(LineString); ok {
		return v.Geometry.Length() < v2.Geometry.Length(), nil
	}
	return LineStringKind < other.Kind(), nil
}

func (v LineString) Hash(nbf *NomsBinFormat) (hash.Hash, error) {
	return getHash(v, nbf)
}

func (v LineString) isPrimitive() bool {
	return true
}

func (v LineString) walkRefs(nbf *NomsBinFormat, cb RefCallback) error {
	return nil
}

func (v LineString) typeOf() (*Type, error) {
	return PrimitiveTypeMap[LineStringKind], nil
}

func (v LineString) Kind() NomsKind {
	return LineStringKind
}

func (v LineString) valueReadWriter() ValueReadWriter {
	return nil
}

func (v LineString) writeTo(w nomsWriter, nbf *NomsBinFormat) error {
	err := LineStringKind.writeTo(w, nbf)
	if err != nil {
		return err
	}

	buf := SerializeLineString(v)
	w.writeString(string(buf))
	return nil
}

func readLineString(nbf *NomsBinFormat, b *valueDecoder) (LineString, error) {
	buf := []byte(b.ReadString())
	srid, _, geomType, err := DeserializeEWKBHeader(buf)
	if err != nil {
		return LineString{}, err
	}
	if geomType != WKBLineID {
		return LineString{}, errors.New("not a linestring")
	}
	buf = buf[EWKBHeaderSize:]
	return DeserializeTypesLine(buf, false, srid), nil
}

func (v LineString) readFrom(nbf *NomsBinFormat, b *binaryNomsReader) (Value, error) {
	buf := []byte(b.ReadString())
	srid, _, geomType, err := DeserializeEWKBHeader(buf)
	if err != nil {
		return LineString{}, err
	}
	if geomType != WKBLineID {
		return LineString{}, errors.New("not a linestring")
	}
	buf = buf[EWKBHeaderSize:]
	return DeserializeTypesLine(buf, false, srid), nil
}

func (v LineString) skip(nbf *NomsBinFormat, b *binaryNomsReader) {
	b.skipString()
}

func (v LineString) HumanReadableString() string {
	coordinateSequence := v.Geometry.CoordSeq()
	points := make([]string, coordinateSequence.Size())

	for i := range coordinateSequence.Size() {
		points[i] = fmt.Sprintf("SRID: %d POINT(%s %s)", v.Geometry.SRID(), strconv.FormatFloat(coordinateSequence.X(i), 'g', -1, 64), strconv.FormatFloat(coordinateSequence.Y(i), 'g', -1, 64))
	}
	s := fmt.Sprintf("SRID: %d LINESTRING(%s)", v.Geometry.SRID(), strings.Join(points, ","))
	return strconv.Quote(s)
}
