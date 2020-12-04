import React from 'react';

import { Button, Input, Overlay } from '@influxdata/clockface';
import { UnControlled as ReactCodeMirror } from 'react-codemirror2';

import 'codemirror/lib/codemirror.css';
import 'codemirror/theme/material.css';

require('codemirror/mode/yaml/yaml');

type Props = {
  onDismiss: () => void;
  value?: string;
  onChange?: (v: string) => void;
};

const OtclOverlay: React.FC<Props> = (props) => {
  const { value = '', onChange, onDismiss } = props;

  const options = {
    tabIndex: 1,
    mode: 'yaml',
    readonly: true,
    lineNumbers: true,
    autoRefresh: true,
    theme: 'material',
    completeSingle: false,
  };

  return (
    <Overlay visible>
      <Overlay.Container maxWidth={800}>
        <Overlay.Header
          title="new"
          onDismiss={() => {
            console.log('dismiss');
          }}
        >
          <Input placeholder="New" />
        </Overlay.Header>

        <Overlay.Body>
          <ReactCodeMirror
            autoCursor
            options={options}
            value={value}
            onChange={(value) => {
              if (onChange !== undefined) {
                onChange(value);
              }
            }}
          />
        </Overlay.Body>

        <Overlay.Footer>
          <Button text="Close" onClick={onDismiss} />
        </Overlay.Footer>
      </Overlay.Container>
    </Overlay>
  );
};

export default OtclOverlay;
