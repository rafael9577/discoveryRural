import { View, Text, TouchableOpacity, ImageBackground } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { useNavigation } from "@react-navigation/native";
import styled from "./styled";

export default () => {
  const navigation = useNavigation();

  const tocard = () => {
    navigation.navigate("infoAnimal");
  };

  const tocard1 = () => {
    navigation.navigate("infoImportantes");
  };

  return (
    <SafeAreaView style={styled.safeArea}>
      {/* Card Principal */}
      <TouchableOpacity style={styled.cardContainer} onPress={tocard} activeOpacity={0.8}>
        <ImageBackground
          source={{ uri: "https://www.infoescola.com/wp-content/uploads/2018/03/sagui-de-tufos-brancos_790108426-1000x667.jpg" }}
          style={styled.cardBackground}
          imageStyle={styled.imageStyle}
        >
          <View style={styled.overlay}>
            <Text style={styled.cardText}>{"Meu Card"}</Text>
          </View>
        </ImageBackground>
      </TouchableOpacity>

      {/* Segundo Card */}
      <TouchableOpacity style={styled.cardContainer2} onPress={tocard1} activeOpacity={0.8}>
        <View style={styled.overlay2}>
          <Text style={styled.cardText}>{"Guia de Emergência"}</Text>
        </View>
      </TouchableOpacity>
    </SafeAreaView>
  );
};