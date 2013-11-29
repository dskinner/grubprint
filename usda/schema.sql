drop table if exists Food;
create table Food (
	Id                 text primary key,
	FoodGroupId        text,
	LongDesc           text,
	ShortDesc          text,
	CommonNames        text,
	ManufacturerName   text,
	Survey             text,
	RefuseDesc         text,
	Refuse             real,
	ScientificName     text,
	NitrogenFactor     real,
	ProteinFactor      real,
	FatFactor          real,
	CarbohydrateFactor real
);

drop table if exists FoodGroup;
create table FoodGroup (
	Id        text primary key,
	Description text
);

drop table if exists LanguaLFactor;
create table LanguaLFactor (
	FoodId          text,
	LanguaLFactorDescriptionId text
);

drop table if exist LanguaLFactorDescription;
create table LanguaLFactorDescription (
	Id          text primary key,
	Description text
);

drop table if exists NutrientData;
create table NutrientData (
	Id               text primary key,
	FoodId           text,
	Value            real,
	DataPoints       real,
	StdError         real,
	SourceCodeId     text,
	DataDerivationId text,
	RefFoodId        text,
	AddNutrMark      text,
	NumStudies       real,
	Min              real,
	Max              real,
	DF               real,
	LowEB            real,
	UpEB             real,
	StatCmt          text,
	AddModDate       text,
	CC               text
);

drop table if exists NutrientDataDefinition;
create table NutrientDataDefinition (
	NutrientDataId text primary key,
	Units          text,
	TagName        text,
	NutrDesc       text,
	NumDec         text,
	Order          real
);

drop table if exists SourceCode;
create table SourceCode (
	Id          text primary key,
	Description text
);

drop table if exists DataDerivation;
create table DataDerivation (
	Id          text primary key,
	Description text
);

drop table if exists Weight;
create table Weight (
	FoodId      text primary key,
	Seq         text,
	Amount      real,
	Description text,
	Grams       real,
	DataPoints  real,
	StdDev      real
);

drop table if exists Footnote;
create table FootNote (
	Id             text,
	FoodId         text,
	Type           text,
	NutrientDataId text,
	Description    text
);

drop table if exists SourcesOfDataLink;
create table SourcesOfDataLink (
	FoodId          text,
	NutrientDataId  text,
	SourcesOfDataId text
);

drop table if exists SourcesOfData;
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
