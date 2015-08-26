-- import usda ascii data
create or replace function public.reduce_dim(anyarray) returns setof anyarray as $$
declare s $1%type;
begin
foreach s slice 1 in array $1
loop
	return next s;
end loop;
end; $$ language plpgsql immutable;

do $$
declare data text[] := array[
	['FoodGroup', 'FD_GROUP'],
	['Food', 'FOOD_DES'],
	['LanguaLFactorDesc', 'LANGDESC'],
	['LanguaLFactor', 'LANGUAL'],
	['SourceCode', 'SRC_CD'],
	['DataDerivation', 'DERIV_CD'],
	['NutrientDef', 'NUTR_DEF'],
	['NutrientData', 'NUT_DATA'],
	['Weight', 'WEIGHT'],
	['Footnote', 'FOOTNOTE'],
	['SourcesOfData', 'DATA_SRC'],
	['SourcesOfDataLink', 'DATSRCLN']
	];
declare r record;
begin
for r in select arr[1] as tablename, arr[2] as filename from reduce_dim(data) vals (arr)
loop
	set client_encoding to 'latin1';
	execute 'copy ' || r.tablename || ' from ''/data/usda/' || r.filename || '.txt'' with delimiter ''^'' csv quote ''~''';
	set client_encoding to default;
end loop;
end $$;
