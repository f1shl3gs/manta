import { useCallback, useEffect, useState } from 'react';
import constate from 'constate';
import { useHistory } from 'react-router-dom';

const escapeKeyCode = 27;

const dispatchResizeEvent = () => {
  setTimeout(() => {
    // Uses longer event object creation method due to IE compatibility.
    const ev = document.createEvent('HTMLEvents')
    ev.initEvent('resize', false, true)
    window.dispatchEvent(ev);
    console.log('dispatched resize event')
  }, 50)
};

const [PresentationModeProvider, usePresentationMode] = constate(
  () => {
    const [inPresentationMode, setInPresentationMode] = useState(false);

    const toggle = useCallback(() => {
      setInPresentationMode(!inPresentationMode);
      dispatchResizeEvent()
    }, [inPresentationMode]);
    const escapePresentationMode = useCallback((event) => {
      if (event.key === 'Escape' || event.keyCode === escapeKeyCode) {
        setInPresentationMode(false);
        dispatchResizeEvent()
      }
    }, []);
    const history = useHistory();

    useEffect(() => {
      window.addEventListener('keyup', escapePresentationMode);
      const unlisten = history.listen(() => {
        setInPresentationMode(false);
        dispatchResizeEvent()
      });

      return () => {
        window.removeEventListener('keyup', escapePresentationMode);
        unlisten();
      };
    });

    return {
      inPresentationMode,
      togglePresentationMode: toggle
    };
  },
  value => value
);

export {
  PresentationModeProvider,
  usePresentationMode
};