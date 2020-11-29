import React from 'react';
import { Overlay } from '@influxdata/clockface';

type Props = {
  title: string;
  onDismiss: () => void;
  children: React.ReactNode;
};

const OtclOverlay: React.FC<Props> = (props) => {
  const { title, children, onDismiss } = props;

  return (
    <Overlay visible={true}>
      <Overlay.Container maxWidth={800}>
        <Overlay.Header title={title} onDismiss={onDismiss} />

        <Overlay.Body>{children}</Overlay.Body>
      </Overlay.Container>
    </Overlay>
  );
};

export default OtclOverlay;
