import React from 'react';
import {
  Button,
  ButtonType,
  Columns,
  ComponentColor,
  ComponentStatus,
  Form,
  Grid,
  Input,
} from '@influxdata/clockface';
import { UnControlled as ReactCodeMirror } from 'react-codemirror2';
import { useOtcl } from 'otcl/state';

import 'codemirror/lib/codemirror.css';
import 'codemirror/theme/material.css';

require('codemirror/mode/yaml/yaml');

const options = {
  tabIndex: 1,
  mode: 'yaml',
  readonly: true,
  lineNumbers: true,
  autoRefresh: true,
  theme: 'material',
  completeSingle: false,
};

const notEmpty = (name: string): string | null => {
  if (!name) {
    return 'Value cannot be empty';
  }

  return null;
};

type Props = {
  onSubmit: () => void;
  onDismiss: () => void;
};

const OtclForm: React.FC<Props> = (props) => {
  const { onSubmit, onDismiss } = props;

  const { otcl, setOtcl } = useOtcl();

  return (
    <Form onSubmit={onSubmit}>
      <Grid>
        <Grid.Row>
          <Grid.Column widthSM={Columns.Six}>
            <Form.ValidationElement
              label={'Name'}
              value={otcl.name || ''}
              validationFunc={notEmpty}
            >
              {(status) => (
                <Input
                  value={otcl.name || ''}
                  name="name"
                  onChange={(ev) =>
                    setOtcl({
                      ...otcl,
                      name: ev.target.value,
                    })
                  }
                  titleText="Name"
                  placeholder="Name this Otcl"
                  autoFocus={true}
                  status={status}
                />
              )}
            </Form.ValidationElement>
          </Grid.Column>

          <Grid.Column widthSM={Columns.Six}>
            <Form.ValidationElement
              label={'Desc'}
              value={otcl.desc || ''}
              validationFunc={notEmpty}
            >
              {(status) => (
                <Input
                  value={otcl.desc || ''}
                  name="desc"
                  onChange={(v) => console.log('desc onchange', v)}
                  titleText="Desc"
                  placeholder="Discribe this Otcl"
                  status={status}
                />
              )}
            </Form.ValidationElement>
          </Grid.Column>
        </Grid.Row>

        <Grid.Row>
          <Form.ValidationElement
            label={'Content'}
            value={otcl.content || ''}
            validationFunc={notEmpty}
          >
            {(status) => (
              <ReactCodeMirror
                autoCursor={true}
                options={options}
                value={otcl.content}
                onChange={(editor, data, value) => {
                  setOtcl({
                    ...otcl,
                    content: value,
                  });
                }}
              />
            )}
          </Form.ValidationElement>
        </Grid.Row>

        <Grid.Row>
          <Grid.Column>
            <Form.Footer>
              <Button
                text="Cancel"
                onClick={onDismiss}
                testID="create-scraper--cancel"
              />
              <Button
                status={ComponentStatus.Default}
                text="Create"
                color={ComponentColor.Success}
                testID="create-scraper--submit"
                type={ButtonType.Submit}
              />
            </Form.Footer>
          </Grid.Column>
        </Grid.Row>
      </Grid>
    </Form>
  );
};

export default OtclForm;
