import { Form, Input, Modal, message } from "antd";
import { useForm } from "antd/es/form/Form";
import TextArea from "antd/es/input/TextArea";
import { friendAdd } from "../../api";

interface AddFriendModalProps {
  isOpen: boolean;
  handleClose: () => void
}

const layout = {
  labelCol: { span: 6 },
  wrapperCol: { span: 18 }
}

export interface AddFriend {
  username: string;
  reason: string;
}

export function AddFriendModal(props: AddFriendModalProps) {

  const [form] = useForm<AddFriend>();

  const handleOk = async function () {
    await form.validateFields();

    const values = form.getFieldsValue();

    try {
      await friendAdd(values);
      message.success('好友申请已发送');
      form.resetFields();
      props.handleClose();
    } catch (e: unknown) {
      if (e instanceof Error) {
        message.error(e.message);
      } else {
        message.error('系统繁忙，请稍后再试');
      }
    }
  }

  return <Modal
    title="添加好友"
    open={props.isOpen}
    onOk={handleOk}
    onCancel={() => props.handleClose()}
    okText={'发送好友请求'}
    cancelText={'取消'}
  >
    <Form
      form={form}
      colon={false}
      {...layout}
    >
      <Form.Item
        label="用户名"
        name="username"
        rules={[
          { required: true, message: '请输入用户名!' },
        ]}
      >
        <Input />
      </Form.Item>
      <Form.Item
        label="添加理由"
        name="reason"
        rules={[
          { required: true, message: '请输入添加理由!' },
        ]}
      >
        <TextArea />
      </Form.Item>
    </Form>
  </Modal>
}
