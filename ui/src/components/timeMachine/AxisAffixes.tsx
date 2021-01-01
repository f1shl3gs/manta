import React, { ChangeEvent, useCallback } from 'react';
import { Columns, FormElement, Grid, Input } from '@influxdata/clockface';

interface Props {
  prefix: string
  suffix: string
  axisName: string
  onSetAxisPrefix: (prefix: string) => void
  onSetAxisSuffix: (suffix: string) => void
}

const AxisAffixes: React.FC<Props> = props => {
  const {
    prefix, suffix, axisName,
    onSetAxisPrefix,
    onSetAxisSuffix
  } = props;

  const onPrefixChange = useCallback((ev: ChangeEvent<HTMLInputElement>) => {
    onSetAxisPrefix(ev.target.value);
  }, []);
  const onSuffixChange = useCallback((ev: ChangeEvent<HTMLInputElement>) => {
    onSetAxisSuffix(ev.target.value);
  }, []);

  return (
    <>
      <Grid.Column widthSM={Columns.Six}>
        <FormElement label={`${axisName.toUpperCase()} Axis Prefix`}>
          <Input value={prefix} onChange={onPrefixChange} />
        </FormElement>
      </Grid.Column>

      <Grid.Column widthSM={Columns.Six}>
        <FormElement label={`${axisName.toUpperCase()} Axis Suffix`}>
          <Input value={suffix} onChange={onSuffixChange} />
        </FormElement>
      </Grid.Column>
    </>
  );
};

export default AxisAffixes;
