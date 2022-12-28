// Libraries
import React, {FunctionComponent} from 'react'
import {useDispatch} from 'react-redux'

// Components
import {
  ButtonShape,
  ComponentColor,
  ComponentSize,
  ConfirmationButton,
  IconFont,
  ResourceCard,
  SquareButton,
} from '@influxdata/clockface'
import CopyToClipboard from 'src/shared/components/CopyToClipboard'

// Types
import {Secret} from 'src/types/secrets'

// Actions
import {deleteSecret} from 'src/secrets/actions/thunks'
import {error, info} from 'src/shared/actions/notifications'

// Utils
import {fromNow} from 'src/shared/utils/duration'
import {push} from '@lagunovsky/redux-react-router'

interface Props {
  secret: Secret
}

const SecretCard: FunctionComponent<Props> = ({secret}) => {
  const dispatch = useDispatch()
  const handleDelete = () => {
    dispatch(deleteSecret(secret.key))
  }
  const handleCopy = (key: string, success: boolean) => {
    if (success) {
      dispatch(info(`Secret key "${secret.key}" has been copied to clipboard`))
    } else {
      dispatch(error(`Copy secret key "${secret.key}" failed`))
    }
  }

  const context = (
    <>
      <CopyToClipboard text={secret.key} onCopy={handleCopy}>
        <SquareButton
          color={ComponentColor.Colorless}
          icon={IconFont.Clipboard_New}
          size={ComponentSize.ExtraSmall}
          shape={ButtonShape.StretchToFit}
          titleText={'Copy to clipboard'}
          testID={'context-copy-to-clipboard--button'}
        />
      </CopyToClipboard>

      <ConfirmationButton
        color={ComponentColor.Colorless}
        icon={IconFont.Trash_New}
        shape={ButtonShape.Square}
        size={ComponentSize.ExtraSmall}
        testID={'context-delete--button'}
        confirmationLabel={'Delete this secret'}
        confirmationButtonText={'Confirm'}
        onConfirm={handleDelete}
      />
    </>
  )

  return (
    <ResourceCard contextMenu={context}>
      <ResourceCard.Name
        name={secret.key}
        onClick={() =>
          dispatch(push(`${window.location.pathname}/${secret.key}/edit`))
        }
      />
      <ResourceCard.Meta>
        <>Modified: {fromNow(secret.updated)}</>
      </ResourceCard.Meta>
    </ResourceCard>
  )
}

export default SecretCard
