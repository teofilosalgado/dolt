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
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/dolthub/dolt/go/store/hash"
	"github.com/twpayne/go-geos"
)

// MultiLineString is a Noms Value wrapper around a string.
type MultiLineString struct {
	Geometry *geos.Geom
}

// Value interface
func (v MultiLineString) Value(ctx context.Context) (Value, error) {
	return v, nil
}

func (v MultiLineString) Equals(other Value) bool {
	if v2, ok := other.(MultiLineString); ok {
		return v.Geometry.Equals(v2.Geometry)
	}
	return false
}

func (v MultiLineString) Less(ctx context.Context, nbf *NomsBinFormat, other LesserValuable) (bool, error) {
	if v2, ok := other.(MultiLineString); ok {
		return v.Geometry.Bounds().Geom().Area() < v2.Geometry.Bounds().Geom().Area(), nil
	}
	return MultiPointKind < other.Kind(), nil
}

func (v MultiLineString) Hash(nbf *NomsBinFormat) (hash.Hash, error) {
	return getHash(v, nbf)
}

func (v MultiLineString) isPrimitive() bool {
	return true
}

func (v MultiLineString) walkRefs(nbf *NomsBinFormat, cb RefCallback) error {
	return nil
}

func (v MultiLineString) typeOf() (*Type, error) {
	return PrimitiveTypeMap[MultiLineStringKind], nil
}

func (v MultiLineString) Kind() NomsKind {
	return MultiLineStringKind
}

func (v MultiLineString) valueReadWriter() ValueReadWriter {
	return nil
}

func (v MultiLineString) writeTo(w nomsWriter, nbf *NomsBinFormat) error {
	err := MultiLineStringKind.writeTo(w, nbf)
	if err != nil {
		return err
	}

	w.writeString(string(SerializeMultiLineString(v)))
	return nil
}

func readMultiLineString(nbf *NomsBinFormat, b *valueDecoder) (MultiLineString, error) {
	buf := []byte(b.ReadString())
	srid, _, geomType, err := DeserializeEWKBHeader(buf)
	if err != nil {
		return MultiLineString{}, err
	}
	if geomType != WKBMultiLineID {
		return MultiLineString{}, errors.New("not a multilinestring")
	}
	buf = buf[EWKBHeaderSize:]
	return DeserializeTypesMLine(buf, false, srid), nil
}

func (v MultiLineString) readFrom(nbf *NomsBinFormat, b *binaryNomsReader) (Value, error) {
	buf := []byte(b.ReadString())
	srid, _, geomType, err := DeserializeEWKBHeader(buf)
	if err != nil {
		return MultiLineString{}, err
	}
	if geomType != WKBMultiLineID {
		return MultiLineString{}, errors.New("not a multilinestring")
	}
	buf = buf[EWKBHeaderSize:]
	return DeserializeTypesMLine(buf, false, srid), nil
}

func (v MultiLineString) skip(nbf *NomsBinFormat, b *binaryNomsReader) {
	b.skipString()
}

func (v MultiLineString) HumanReadableString() string {
	coordinateSequence := v.Geometry.CoordSeq()
	points := make([]string, coordinateSequence.Size())

	for i := range coordinateSequence.Size() {
		points[i] = fmt.Sprintf("SRID: %d POINT(%s %s)", v.Geometry.SRID(), strconv.FormatFloat(coordinateSequence.X(i), 'g', -1, 64), strconv.FormatFloat(coordinateSequence.Y(i), 'g', -1, 64))
	}
	s := fmt.Sprintf("SRID: %d MULTILINESTRING(%s)", v.Geometry.SRID(), strings.Join(points, ","))
	return strconv.Quote(s)
}
