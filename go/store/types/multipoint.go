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

// MultiPoint is a Noms Value wrapper around a string.
type MultiPoint struct {
	Geometry *geos.Geom
}

// Value interface
func (v MultiPoint) Value(ctx context.Context) (Value, error) {
	return v, nil
}

func (v MultiPoint) Equals(other Value) bool {
	if v2, ok := other.(MultiPoint); ok {
		return v.Geometry.Equals(v2.Geometry)
	}
	return false
}

func (v MultiPoint) Less(ctx context.Context, nbf *NomsBinFormat, other LesserValuable) (bool, error) {
	if v2, ok := other.(MultiPoint); ok {
		return v.Geometry.Bounds().Geom().Area() < v2.Geometry.Bounds().Geom().Area(), nil
	}
	return MultiPointKind < other.Kind(), nil
}

func (v MultiPoint) Hash(nbf *NomsBinFormat) (hash.Hash, error) {
	return getHash(v, nbf)
}

func (v MultiPoint) isPrimitive() bool {
	return true
}

func (v MultiPoint) walkRefs(nbf *NomsBinFormat, cb RefCallback) error {
	return nil
}

func (v MultiPoint) typeOf() (*Type, error) {
	return PrimitiveTypeMap[MultiPointKind], nil
}

func (v MultiPoint) Kind() NomsKind {
	return MultiPointKind
}

func (v MultiPoint) valueReadWriter() ValueReadWriter {
	return nil
}

func (v MultiPoint) writeTo(w nomsWriter, nbf *NomsBinFormat) error {
	err := MultiPointKind.writeTo(w, nbf)
	if err != nil {
		return err
	}
	w.writeString(string(SerializeMultiPoint(v)))
	return nil
}

func readMultiPoint(nbf *NomsBinFormat, b *valueDecoder) (MultiPoint, error) {
	buf := []byte(b.ReadString())
	srid, _, geomType, err := DeserializeEWKBHeader(buf)
	if err != nil {
		return MultiPoint{}, nil
	}
	if geomType != WKBMultiPointID {
		return MultiPoint{}, errors.New("not a multipoint")
	}
	buf = buf[EWKBHeaderSize:]
	return DeserializeTypesMPoint(buf, false, srid), nil
}

func (v MultiPoint) readFrom(nbf *NomsBinFormat, b *binaryNomsReader) (Value, error) {
	buf := []byte(b.ReadString())
	srid, _, geomType, err := DeserializeEWKBHeader(buf)
	if err != nil {
		return MultiPoint{}, nil
	}
	if geomType != WKBMultiPointID {
		return MultiPoint{}, errors.New("not a multipoint")
	}
	return DeserializeTypesMPoint(buf, false, srid), nil
}

func (v MultiPoint) skip(nbf *NomsBinFormat, b *binaryNomsReader) {
	b.skipString()
}

func (v MultiPoint) HumanReadableString() string {
	coordinateSequence := v.Geometry.CoordSeq()
	points := make([]string, coordinateSequence.Size())

	for i := range coordinateSequence.Size() {
		points[i] = fmt.Sprintf("SRID: %d POINT(%s %s)", v.Geometry.SRID(), strconv.FormatFloat(coordinateSequence.X(i), 'g', -1, 64), strconv.FormatFloat(coordinateSequence.Y(i), 'g', -1, 64))
	}
	s := fmt.Sprintf("SRID: %d MULTIPOINT(%s)", v.Geometry.SRID(), strings.Join(points, ","))
	return strconv.Quote(s)
}
