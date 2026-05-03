const API_BASE_URL = "http://localhost:8081";
const API_BASE_URL_SECURE = "https://localhost:8081";
const API_BASE_URL_NGROK = "https://27dc-213-176-17-134.ngrok-free.app";

const routes = {
  login: `${API_BASE_URL}/v1/login`,
  loginSecure: `${API_BASE_URL_SECURE}/v1/login`,
  loginNgrok: `${API_BASE_URL_NGROK}/v1/login`,

  getAllDefaultParsingChannelsWithCategories: `${API_BASE_URL}/v1/get_all_default_parsing_channels_with_categories`,
  addNewDefaultParsingChannel: `${API_BASE_URL}/v1/add_new_default_parsing_channel`,
  deleteParsingChannel: `${API_BASE_URL}/v1/delete_parsing_channel`,

  getRecommendatedPosts: `${API_BASE_URL}/v1/get_recommendated_posts`,

  getUserPriorityChannels: `${API_BASE_URL}/v1/get_user_priority_channels`,
  setPriorityChannels: `${API_BASE_URL}/v1/set_priority_channels`,
  getAllDefaultParsingChannels: `${API_BASE_URL}/v1/get_all_default_parsing_channels`,
  addNewUserCustomParsingChannel: `${API_BASE_URL}/v1/add_new_user_custom_parsing_channel`,
  getAllUserCustomParsingChannels: `${API_BASE_URL}/v1/get_all_user_custom_parsing_channels`,
};

export { routes, API_BASE_URL, API_BASE_URL_SECURE, API_BASE_URL_NGROK };
