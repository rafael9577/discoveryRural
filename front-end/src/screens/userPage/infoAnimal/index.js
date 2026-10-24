import { View, Text, TouchableOpacity, ImageBackground } from "react-native";
import { SafeAreaView, } from "react-native-safe-area-context";

import styled from "./styled";

export default () => {

    return (
        <SafeAreaView style={styled.safeArea}>
            <Text style={{alignSelf:'center', marginTop:50, fontSize: 25}}>Informações do animal</Text>
            <View style={{width:"75%", height:"75%", alignSelf: "center"}}>
                <ImageBackground source={{ uri: "https://www.infoescola.com/wp-content/uploads/2018/03/sagui-de-tufos-brancos_790108426-1000x667.jpg" }} style={styled.cardBackground} imageStyle={styled.imageStyle} />
            </View>
        </SafeAreaView>
)};