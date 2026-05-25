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

// MultiPolygon is a Noms Value wrapper around a string.
type MultiPolygon struct {
	Geometry *geos.Geom
}

// Value interface
func (v MultiPolygon) Value(ctx context.Context) (Value, error) {
	return v, nil
}

func (v MultiPolygon) Equals(other Value) bool {
	if v2, ok := other.(MultiPolygon); ok {
		return v.Geometry.Equals(v2.Geometry)
	}
	return false
}

func (v MultiPolygon) Less(ctx context.Context, nbf *NomsBinFormat, other LesserValuable) (bool, error) {
	if v2, ok := other.(MultiPolygon); ok {
		return v.Geometry.Bounds().Geom().Area() < v2.Geometry.Bounds().Geom().Area(), nil
	}
	return MultiPointKind < other.Kind(), nil
}

func (v MultiPolygon) Hash(nbf *NomsBinFormat) (hash.Hash, error) {
	return getHash(v, nbf)
}

func (v MultiPolygon) isPrimitive() bool {
	return true
}

func (v MultiPolygon) walkRefs(nbf *NomsBinFormat, cb RefCallback) error {
	return nil
}

func (v MultiPolygon) typeOf() (*Type, error) {
	return PrimitiveTypeMap[MultiPolygonKind], nil
}

func (v MultiPolygon) Kind() NomsKind {
	return MultiPolygonKind
}

func (v MultiPolygon) valueReadWriter() ValueReadWriter {
	return nil
}

func (v MultiPolygon) writeTo(w nomsWriter, nbf *NomsBinFormat) error {
	err := MultiPolygonKind.writeTo(w, nbf)
	if err != nil {
		return err
	}

	w.writeString(string(SerializeMultiPolygon(v)))
	return nil
}

func readMultiPolygon(nbf *NomsBinFormat, b *valueDecoder) (MultiPolygon, error) {
	buf := []byte(b.ReadString())
	srid, _, geomType, err := DeserializeEWKBHeader(buf)
	if err != nil {
		return MultiPolygon{}, err
	}
	if geomType != WKBMultiPolyID {
		return MultiPolygon{}, errors.New("not a multipolygon")
	}
	buf = buf[EWKBHeaderSize:]
	return DeserializeTypesMPoly(buf, false, srid), nil
}

func (v MultiPolygon) readFrom(nbf *NomsBinFormat, b *binaryNomsReader) (Value, error) {
	buf := []byte(b.ReadString())
	srid, _, geomType, err := DeserializeEWKBHeader(buf)
	if err != nil {
		return MultiPolygon{}, err
	}
	if geomType != WKBMultiPolyID {
		return MultiPolygon{}, errors.New("not a multipolygon")
	}
	buf = buf[EWKBHeaderSize:]
	return DeserializeTypesMPoly(buf, false, srid), nil
}

func (v MultiPolygon) skip(nbf *NomsBinFormat, b *binaryNomsReader) {
	b.skipString()
}

func (v MultiPolygon) HumanReadableString() string {
	coordinateSequence := v.Geometry.CoordSeq()
	points := make([]string, coordinateSequence.Size())

	for i := range coordinateSequence.Size() {
		points[i] = fmt.Sprintf("SRID: %d POINT(%s %s)", v.Geometry.SRID(), strconv.FormatFloat(coordinateSequence.X(i), 'g', -1, 64), strconv.FormatFloat(coordinateSequence.Y(i), 'g', -1, 64))
	}
	s := fmt.Sprintf("SRID: %d MULTIPOLYGON(%s)", v.Geometry.SRID(), strings.Join(points, ","))
	return strconv.Quote(s)
}
