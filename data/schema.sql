-- TODO(d)
-- select * from food where longdesc ~* 'navy|beans';
-- CREATE EXTENSION pg_trgm;
-- create index food_longdesc_idx on food using gin(longdesc gin_trgm_ops);

drop table if exists FoodGroup cascade;
create table FoodGroup (
	Id          text primary key,
	Description text
);

drop table if exists Food cascade;
create table Food (
	Id                 text primary key,
	FoodGroupId        text references FoodGroup (Id),
	LongDesc           text,
	ShortDesc          text,
	CommonNames        text,
	Manufacturer       text,
	Survey             text,
	RefuseDesc         text,
	Refuse             real,
	ScientificName     text,
	NitrogenFactor     real,
	ProteinFactor      real,
	FatFactor          real,
	CarbohydrateFactor real
);

drop table if exists LanguaLFactorDesc cascade;
create table LanguaLFactorDesc (
	Id          text primary key,
	Description text
);

drop table if exists LanguaLFactor cascade;
create table LanguaLFactor (
	FoodId              text references Food (Id),
	LanguaLFactorDescId text references LanguaLFactorDesc (Id),
	primary key (FoodId, LanguaLFactorDescId)
);

drop table if exists SourceCode cascade;
create table SourceCode (
	Id          text primary key,
	Description text
);

drop table if exists DataDerivation cascade;
create table DataDerivation (
	Id          text primary key,
	Description text
);

drop table if exists NutrientDef cascade;
create table NutrientDef (
	Id       text primary key,
	Units    text,
	TagName  text,
	NutrDesc text,
	NumDec   text,
	Sort     real
);

drop table if exists NutrientData cascade;
create table NutrientData (
	FoodId           text references Food (Id),
	NutrientDefId    text references NutrientDef (Id),
	Value            real,
	DataPoints       real,
	StdError         real,
	SourceCodeId     text references SourceCode (Id),

	-- TODO(d) allow nullable but copy from not setting null from empty string?
	DataDerivationId text, -- references DataDerivation (Id),
	RefFoodId        text, -- references Food (Id),

	AddNutrMark      text,
	NumStudies       real,
	Min              real,
	Max              real,
	DF               real,
	LowEB            real,
	UpEB             real,
	StatCmt          text,
	AddModDate       text,
	CC               text,
	primary key (FoodId, NutrientDefId)
);

drop table if exists Weight cascade;
create table Weight (
	FoodId      text references Food (Id),
	Seq         text,
	Amount      real,
	Description text,
	Grams       real,
	DataPoints  real,
	StdDev      real,
	primary key (FoodId, Seq)
);

drop table if exists Footnote cascade;
create table Footnote (
	FoodId        text references Food (Id),
	Seq           text,
	Type          text,
	NutrientDefId text, -- references NutrientDef (Id),
	Description   text
);

drop table if exists SourcesOfData cascade;
create table SourcesOfData (
	Id         text primary key,
	Authors    text,
	Title      text,
	Year       text,
	Journal    text,
	VolCity    text,
	IssueState text,
	StartPage  text,
	EndPage    text
);

drop table if exists SourcesOfDataLink;
create table SourcesOfDataLink (
	FoodId          text references Food (Id),
	NutrientDefId   text references NutrientDef (Id),
	SourcesOfDataId text references SourcesOfData (Id),
	primary key (FoodId, NutrientDefId, SourcesOfDataId)
);

