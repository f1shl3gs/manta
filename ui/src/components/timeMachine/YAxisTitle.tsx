import React, { useCallback } from 'react';
import { Columns, Form, Grid, Input } from '@influxdata/clockface';

interface Props {
  label: string
  onUpdateYAxisLabel: (label: string) => void
}

const YAxisTitle: React.FC<Props> = props => {
  const { label, onUpdateYAxisLabel } = props;
  const onChange = useCallback((ev: React.ChangeEvent<HTMLInputElement>) => {
    onUpdateYAxisLabel(ev.target.value);
  }, []);

  return (
    <Grid.Column widthXS={Columns.Twelve}>
      <Form.Element label={'Y Axis Label'}>
        <Input value={label} onChange={onChange} />
      </Form.Element>
    </Grid.Column>
  );
};

export default YAxisTitle;