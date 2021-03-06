import React from 'react';
import PropTypes from 'prop-types';
import { Menu, Dropdown, Icon, Button } from 'antd';

class EllipsisMenu extends React.Component {
  static propTypes = {
    placement: PropTypes.oneOf([
      'bottomLeft',
      'bottomCenter',
      'bottomRight',
      'topLeft',
      'topCenter',
      'topRight',
    ]),
    menuText: PropTypes.oneOfType([PropTypes.string, PropTypes.array]),
    menuFunc: PropTypes.oneOfType([PropTypes.func, PropTypes.array]),
    disabled: PropTypes.oneOfType([PropTypes.bool, PropTypes.array]),
    disablAll: PropTypes.bool,
  };
  static defaultProps = {
    placement: 'bottomLeft',
    disabled: false,
  };

  render() {
    const { placement, menuText, menuFunc, disabled, disablAll } = this.props;
    const _menuText = menuText || intl.get('operation.delete');
    const menu = (
      <Menu>
        {_.isArray(menuText) ? (
          menuText.map((m, i) => (
            <Menu.Item
              key={m}
              onClick={e => {
                e.domEvent.preventDefault();
                e.domEvent.stopPropagation();
                menuFunc[i]();
              }}
              disabled={disabled[i]}
            >
              {m}
            </Menu.Item>
          ))
        ) : (
          <Menu.Item
            onClick={e => {
              e.domEvent.preventDefault();
              e.domEvent.stopPropagation();
              menuFunc();
            }}
            disabled={disabled}
          >
            {_menuText}
          </Menu.Item>
        )}
      </Menu>
    );
    return (
      <div style={{ marginLeft: '8px', display: 'inline-block' }}>
        <Dropdown
          disabled={disablAll}
          overlay={menu}
          placement={placement}
          trigger={['click']}
          onClick={e => {
            e.preventDefault();
            e.stopPropagation();
          }}
        >
          <Button size="small">
            <Icon
              type="ellipsis"
              style={{
                transform: 'rotate(90deg)',
              }}
            />
          </Button>
        </Dropdown>
      </div>
    );
  }
}

export default EllipsisMenu;
