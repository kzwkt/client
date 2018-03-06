// @flow
import * as React from 'react'
import get from 'lodash/get'
import type {Props, ItemProps, TabBarButtonProps} from './tab-bar'
import {NativeTouchableWithoutFeedback, NativeStyleSheet} from './native-wrappers.native'
import Badge from './badge'
import Box from './box'
import Icon from './icon'
import Text from './text'
import {globalStyles, globalColors, globalMargins, styleSheetCreate} from '../styles'

class TabBarItem extends React.Component<ItemProps> {
  render() {
    return this.props.children
  }
}

class SimpleTabBarButton extends React.Component<ItemProps> {
  render() {
    const selectedColor = this.props.selectedColor || globalColors.blue
    return (
      <Box style={[styles.tab, this.props.style]}>
        <Text
          type="BodySmallSemibold"
          style={[styles.label, {color: this.props.selected ? globalColors.black_75 : globalColors.black_40}]}
        >
          {!!this.props.label && this.props.label.toUpperCase()}
        </Text>
        <Box style={this.props.selected ? stylesSelectedUnderline(selectedColor) : styles.unselected} />
      </Box>
    )
  }
}

const UnderlineHighlight = () => (
  <Box
    style={{
      position: 'absolute',
      bottom: 0,
      left: 24,
      right: 24,
      height: 2,
      borderTopLeftRadius: 3,
      borderTopRightRadius: 3,
      backgroundColor: globalColors.white,
    }}
  />
)

const TabBarButton = (props: TabBarButtonProps) => {
  const badgeNumber = props.badgeNumber || 0

  let badgeComponent = null
  if (badgeNumber > 0) {
    if (props.badgePosition === 'top-right') {
      badgeComponent = (
        <Box
          style={{
            ...globalStyles.flexBoxColumn,
            ...globalStyles.fillAbsolute,
            alignItems: 'center',
            justifyContent: 'center',
          }}
        >
          <Badge badgeNumber={badgeNumber} badgeStyle={{marginRight: -30, marginTop: -20}} />
        </Box>
      )
    } else {
      badgeComponent = <Badge badgeNumber={badgeNumber} badgeStyle={{marginLeft: 5}} />
    }
  }

  const content = (
    <Box style={[styles.tabBarButtonIcon, props.style]}>
      <Icon
        type={
          // $FlowIssue
          props.source.icon
        }
        style={[props.isNav ? null : {width: 32}, props.styleIcon]}
      />
      {!!props.label && (
        <Text type="BodySemibold" style={{textAlign: 'center', ...props.styleLabel}}>
          {props.label}
        </Text>
      )}
      {badgeComponent}
      {props.underlined && <UnderlineHighlight />}
    </Box>
  )
  if (props.onClick) {
    return (
      <NativeTouchableWithoutFeedback onPress={props.onClick} style={{flex: 1}}>
        {content}
      </NativeTouchableWithoutFeedback>
    )
  }
  return content
}

class TabBar extends React.Component<Props> {
  _labels(): Array<React.Node> {
    // TODO: Not sure why I have to wrap the child in a box, but otherwise touches won't work
    // $FlowIssue dunno
    return (this.props.children || []).map((item: {props: ItemProps}, i) => {
      const key = item.props.label || get(item, 'props.tabBarButton.props.label') || i
      return (
        <NativeTouchableWithoutFeedback key={key} onPress={item.props.onClick || (() => {})}>
          <Box style={{flex: 1}}>
            <Box style={item.props.styleContainer}>
              {item.props.tabBarButton || <SimpleTabBarButton {...item.props} />}
            </Box>
          </Box>
        </NativeTouchableWithoutFeedback>
      )
    })
  }

  _content(): any {
    return (this.props.children || []).find(i => i.props.selected)
  }

  render() {
    const tabBarButtons = (
      <Box style={globalStyles.flexBoxColumn}>
        <Box style={[globalStyles.flexBoxRow, this.props.styleTabBar]}>{this._labels()}</Box>
        {this.props.underlined && <Box style={styles.underline} />}
      </Box>
    )
    return (
      <Box style={[styles.container, this.props.style]}>
        {!this.props.tabBarOnBottom && tabBarButtons}
        {this._content()}
        {this.props.tabBarOnBottom && tabBarButtons}
      </Box>
    )
  }
}

const stylesSelectedUnderline = color => ({
  height: 3,
  marginBottom: -1,
  alignSelf: 'stretch',
  backgroundColor: color,
})

const styles = styleSheetCreate({
  container: {
    ...globalStyles.flexBoxColumn,
    ...globalStyles.fullHeight,
  },
  underline: {
    height: NativeStyleSheet.hairlineWidth,
    alignSelf: 'stretch',
    backgroundColor: globalColors.black_05,
  },
  unselected: {
    height: 2,
  },
  label: {
    marginTop: 11,
    marginBottom: 11,
    height: globalMargins.small,
  },
  tabBarButtonIcon: {
    ...globalStyles.flexBoxColumn,
    alignItems: 'center',
    flexGrow: 1,
    justifyContent: 'center',
    position: 'relative',
  },
  tab: {
    ...globalStyles.flexBoxColumn,
    alignItems: 'center',
    flexGrow: 1,
    justifyContent: 'flex-end',
  },
})

export {TabBarItem, TabBarButton}

export default TabBar
