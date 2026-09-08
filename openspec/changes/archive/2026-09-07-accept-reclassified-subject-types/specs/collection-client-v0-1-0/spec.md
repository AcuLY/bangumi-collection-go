## MODIFIED Requirements

### Requirement: The public DTO preserves the complete collection record
`Fetch` and `FetchPage` SHALL return `*Subject` values with public fields `ID`, `SubjectID`, `SubjectType`, `Type`, `Name`, `NameCn`, `Rate`, `Comment`, `Tags`, `UpdatedAt`, `VolStatus`, `EpStatus`, and `Private`. `SubjectID` SHALL contain upstream top-level `subject_id`; compatibility field `ID` SHALL equal `SubjectID`. `SubjectType` SHALL contain upstream `subject_type`; `Type` SHALL contain the collection state. Present non-null comment, required tags, update timestamp, progress, rating, and private marker SHALL be preserved without semantic inference.

The decoder SHALL follow the current official Bangumi OAS: top-level
`subject_id`, `subject_type`, `rate`, `type`, `tags`, `ep_status`,
`vol_status`, `updated_at`, and `private` are required, while `comment` and
nested `subject` are optional. An absent or JSON-null comment SHALL map to
`Comment == ""`. An absent nested subject SHALL still set `ID == SubjectID`
and leave `Name`/`NameCn` empty; when present, its required
identity/type/name tuple SHALL be complete. Its ID SHALL match top-level
subject_id, and its type SHALL be a supported subject type. A differing valid
nested type SHALL be accepted without changing the returned SubjectType,
which SHALL continue to preserve top-level subject_type.
A missing/null required tags array SHALL fail, while an empty required array
SHALL become a non-nil empty slice and every returned tag slice SHALL be
copied.

The decoder SHALL require a positive subject ID, valid subject type in `{1,2,3,4,6}`, valid collection type in `{1,2,3,4,5}`, rate in `0..10`, non-negative progress, RFC3339 `updated_at`, and exact normalized requested offset/limit metadata. Every page total SHALL be in `0..1_000_000` per collection type and `len(data)` SHALL NOT exceed the normalized requested limit. Any violation SHALL return `*ProtocolError` matching `ErrProtocol` before allocation or scheduling. Unknown additive upstream JSON fields MAY be ignored; missing/malformed required content, malformed present non-null optional content, multiple JSON values, or trailing non-whitespace SHALL fail.

#### Scenario: Complete valid collection decodes
- **WHEN** an `httptest` response contains every required collection field and a consistent optional nested subject
- **THEN** both page and aggregate operations SHALL return every exact value
- **AND** `ID` SHALL equal `SubjectID`
- **AND** mutating a source decode buffer or another result SHALL NOT mutate the returned tags

#### Scenario: Required field or identity is invalid
- **WHEN** a success payload omits a required field, uses an invalid enum/range/timestamp, disagrees on subject identity, has invalid page metadata, or contains trailing JSON
- **THEN** the operation SHALL return a typed non-retryable decode or protocol error
- **AND** it SHALL return no partial aggregate result

#### Scenario: Optional comment and subject are absent
- **WHEN** a success payload contains every required top-level field but omits `comment` and nested `subject`
- **THEN** decoding SHALL succeed with empty `Comment`, `Name`, and `NameCn`
- **AND** `ID` SHALL still equal the required top-level `SubjectID`

#### Scenario: Optional comment is null
- **WHEN** an otherwise valid success payload contains `comment: null`
- **THEN** page and aggregate decoding SHALL succeed with `Comment == ""`
- **AND** non-null non-string comments SHALL remain typed non-retryable protocol failures

#### Scenario: Empty public collection succeeds
- **WHEN** the first valid page reports total zero and contains no data
- **THEN** `Fetch` SHALL return a non-nil empty slice and nil error

#### Scenario: Nested subject metadata has a different supported type
- **WHEN** an anime collection record has top-level subject_type 2 and a same-ID nested subject with supported type 6
- **THEN** FetchPage and Fetch SHALL retain the record and its names with SubjectType 2
- **AND** the aggregate SHALL not drop the record, rewrite its type or fail because of that supported nested type difference

#### Scenario: Nested metadata remains invalid
- **WHEN** the nested subject has a missing, null, malformed or unsupported type, missing names, or a different ID
- **THEN** the decoder SHALL still return a non-retryable protocol error without partial aggregate data