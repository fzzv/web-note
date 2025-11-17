import { Form, Input, Modal, message } from "antd";
import { useForm } from "antd/es/form/Form";
import { addMember } from "../../api";

interface AddMemberModalProps {
  chatroomId: number;
  isOpen: boolean;
  handleClose: () => void
}

const layout = {
  labelCol: { span: 6 },
  wrapperCol: { span: 18 }
}

export interface AddMember {
  username: string;
}

export function AddMemberModal(props: AddMemberModalProps) {

  const [form] = useForm<AddMember>();

  const handleOk = async function () {
    await form.validateFields();

    const values = form.getFieldsValue();

    try {
      await addMember(props.chatroomId, values.username);

      message.success('成员添加成功');
      form.resetFields();
      props.handleClose();
    } catch (e: unknown) {
      if (e instanceof Error) {
        console.error(e.message);
      } else {
        console.error('Unknown error:', e);
      }
    }
  }

  return <Modal
    title="添加成员"
    open={props.isOpen}
    onOk={handleOk}
    onCancel={() => props.handleClose()}
    okText={'添加'}
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
    </Form>
  </Modal>
}
